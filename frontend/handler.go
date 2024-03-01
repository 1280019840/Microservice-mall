package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"frontend/money"
	pb "frontend/proto"
)

var log *logrus.Logger

func initializeLogger() {
	log = logrus.New()
	log.Level = logrus.DebugLevel
	log.Formatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	}
	log.Out = os.Stdout
}

// 初始化日志
func init() {
	initializeLogger()
}

// 主页
func (fe *FrontendServer) HomeHandler(ctx *gin.Context) {
	r := ctx.Request
	currencies, err := fe.getCurrencies(r.Context())
	if err != nil {
		renderHTTPError(log, ctx, errors.Wrap(err, "不能查询到货币"), http.StatusInternalServerError)
		return
	}
	products, err := fe.getProducts(r.Context())

	if err != nil {
		renderHTTPError(log, ctx, errors.Wrap(err, "不能查询到商品"), http.StatusInternalServerError)
		return
	}
	cart, err := fe.getCart(r.Context(), sessionID(r))
	if err != nil {
		renderHTTPError(log, ctx, errors.Wrap(err, "不能查询到购物车"), http.StatusInternalServerError)
		return
	}

	type productView struct {
		Item  *pb.Product
		Price *pb.Money
	}
	ps := make([]productView, len(products))
	for i, p := range products {
		price, err := fe.convertCurrency(r.Context(), p.GetPriceUsd(), currentCurrency(r))
		if err != nil {
			renderHTTPError(log, ctx, errors.Wrapf(err, "货币转换失败 %s", p.GetId()), http.StatusInternalServerError)
			return
		}
		ps[i] = productView{p, price}
	}

	resultMap := map[string]interface{}{
		"session_id":    sessionID(r),
		"request_id":    r.Context().Value(ctxKeyRequestID{}),
		"user_currency": currentCurrency(r),
		"show_currency": true,
		"currencies":    currencies,
		"products":      ps,
		"cart_size":     cartSize(cart),
		"ad":            fe.chooseAd(r.Context(), []string{}, log),
	}

	ctx.HTML(http.StatusOK, "home", resultMap)

}

// 商品
func (fe *FrontendServer) ProductHandler(ctx *gin.Context) {
	r := ctx.Request
	id := ctx.Param("id")
	if id == "" {
		renderHTTPError(log, ctx, errors.New("商品id没有指定"), http.StatusBadRequest)
		return
	}
	log.WithField("id", id).WithField("currency", currentCurrency(r)).
		Debug("商品服务")

	p, err := fe.getProduct(r.Context(), id)
	if err != nil {
		renderHTTPError(log, ctx, errors.Wrap(err, "不能查询到商品"), http.StatusInternalServerError)
		return
	}
	currencies, err := fe.getCurrencies(r.Context())
	if err != nil {
		renderHTTPError(log, ctx, errors.Wrap(err, "不能查询到货币"), http.StatusInternalServerError)
		return
	}

	cart, err := fe.getCart(r.Context(), sessionID(r))
	if err != nil {
		renderHTTPError(log, ctx, errors.Wrap(err, "不能查询到购物车"), http.StatusInternalServerError)
		return
	}

	price, err := fe.convertCurrency(r.Context(), p.GetPriceUsd(), currentCurrency(r))
	if err != nil {
		renderHTTPError(log, ctx, errors.Wrap(err, "货币转换失败"), http.StatusInternalServerError)
		return
	}

	recommendations, err := fe.getRecommendations(r.Context(), sessionID(r), []string{id})
	if err != nil {
		renderHTTPError(log, ctx, errors.Wrap(err, "商品推荐失败"), http.StatusInternalServerError)
		return
	}

	product := struct {
		Item  *pb.Product
		Price *pb.Money
	}{p, price}

	resultMap := map[string]interface{}{
		"session_id":      sessionID(r),
		"request_id":      r.Context().Value(ctxKeyRequestID{}),
		"ad":              fe.chooseAd(r.Context(), p.Categories, log),
		"user_currency":   currentCurrency(r),
		"show_currency":   true,
		"currencies":      currencies,
		"product":         product,
		"recommendations": recommendations,
		"cart_size":       cartSize(cart),
	}

	ctx.HTML(http.StatusOK, "product", resultMap)
}

// 添加购物车
func (fe *FrontendServer) addToCartHandler(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	quantity, _ := strconv.ParseUint(r.FormValue("quantity"), 10, 32)
	productID := r.FormValue("product_id")
	if productID == "" || quantity == 0 {
		renderHTTPError(log, ctx, errors.New("无效表单输入"), http.StatusBadRequest)
		return
	}
	log.WithField("product", productID).WithField("quantity", quantity).Debug("添加到购物车")

	p, err := fe.getProduct(r.Context(), productID)
	if err != nil {
		renderHTTPError(log, ctx, errors.Wrap(err, "不能查询商品"), http.StatusInternalServerError)
		return
	}

	if err := fe.insertCart(r.Context(), sessionID(r), p.GetId(), int32(quantity)); err != nil {
		renderHTTPError(log, ctx, errors.Wrap(err, "条件购物车失败"), http.StatusInternalServerError)
		return
	}
	w.Header().Set("location", "/cart")
	w.WriteHeader(http.StatusFound)
}

// 清空购物车
func (fe *FrontendServer) emptyCartHandler(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer

	log.Debug("清空购物车")

	if err := fe.emptyCart(r.Context(), sessionID(r)); err != nil {
		renderHTTPError(log, ctx, errors.Wrap(err, "清空购物车失败"), http.StatusInternalServerError)
		return
	}
	w.Header().Set("location", "/")
	w.WriteHeader(http.StatusFound)
}

// 浏览购物车
func (fe *FrontendServer) viewCartHandler(ctx *gin.Context) {
	r := ctx.Request
	log.Debug("浏览购物车")
	currencies, err := fe.getCurrencies(r.Context())
	if err != nil {
		renderHTTPError(log, ctx, errors.Wrap(err, "不能查询到货币"), http.StatusInternalServerError)
		return
	}
	cart, err := fe.getCart(r.Context(), sessionID(r))
	if err != nil {
		renderHTTPError(log, ctx, errors.Wrap(err, "不能查询到购物车"), http.StatusInternalServerError)
		return
	}

	recommendations, err := fe.getRecommendations(r.Context(), sessionID(r), cartIDs(cart))
	if err != nil {
		renderHTTPError(log, ctx, errors.Wrap(err, "获得商品推荐失败"), http.StatusInternalServerError)
		return
	}

	shippingCost, err := fe.getShippingQuote(r.Context(), cart, currentCurrency(r))
	if err != nil {
		renderHTTPError(log, ctx, errors.Wrap(err, "配送失败"), http.StatusInternalServerError)
		return
	}

	type cartItemView struct {
		Item     *pb.Product
		Quantity int32
		Price    *pb.Money
	}
	items := make([]cartItemView, len(cart))
	totalPrice := &pb.Money{CurrencyCode: currentCurrency(r)}
	for i, item := range cart {
		p, err := fe.getProduct(r.Context(), item.GetProductId())
		if err != nil {
			renderHTTPError(log, ctx, errors.Wrapf(err, "不能查询到商品 #%s", item.GetProductId()), http.StatusInternalServerError)
			return
		}
		price, err := fe.convertCurrency(r.Context(), p.GetPriceUsd(), currentCurrency(r))
		if err != nil {
			renderHTTPError(log, ctx, errors.Wrapf(err, "不能转换货币 #%s", item.GetProductId()), http.StatusInternalServerError)
			return
		}

		multPrice := money.MultiplySlow(price, uint32(item.GetQuantity()))
		items[i] = cartItemView{
			Item:     p,
			Quantity: item.GetQuantity(),
			Price:    multPrice,
		}
		totalPrice = money.Must(money.Sum(totalPrice, multPrice))
	}
	totalPrice = money.Must(money.Sum(totalPrice, shippingCost))
	year := time.Now().Year()

	resultMap := map[string]interface{}{
		"session_id":       sessionID(r),
		"request_id":       r.Context().Value(ctxKeyRequestID{}),
		"user_currency":    currentCurrency(r),
		"currencies":       currencies,
		"recommendations":  recommendations,
		"cart_size":        cartSize(cart),
		"shipping_cost":    shippingCost,
		"show_currency":    true,
		"total_cost":       totalPrice,
		"items":            items,
		"expiration_years": []int{year, year + 1, year + 2, year + 3, year + 4},
	}

	ctx.HTML(http.StatusOK, "cart", resultMap)

}

// 下订单
func (fe *FrontendServer) placeOrderHandler(ctx *gin.Context) {
	r := ctx.Request
	log.Debug("下订单")

	var (
		email         = r.FormValue("email")
		streetAddress = r.FormValue("street_address")
		zipCode, _    = strconv.ParseInt(r.FormValue("zip_code"), 10, 32)
		city          = r.FormValue("city")
		state         = r.FormValue("state")
		country       = r.FormValue("country")
		ccNumber      = strings.ReplaceAll(r.FormValue("credit_card_number"), "-", "")
		ccMonth, _    = strconv.ParseInt(r.FormValue("credit_card_expiration_month"), 10, 32)
		ccYear, _     = strconv.ParseInt(r.FormValue("credit_card_expiration_year"), 10, 32)
		ccCVV, _      = strconv.ParseInt(r.FormValue("credit_card_cvv"), 10, 32)
	)

	order, err := fe.checkoutService.PlaceOrder(r.Context(), &pb.PlaceOrderRequest{
		Email: email,
		CreditCard: &pb.CreditCardInfo{
			CreditCardNumber:          ccNumber,
			CreditCardExpirationMonth: int32(ccMonth),
			CreditCardExpirationYear:  int32(ccYear),
			CreditCardCvv:             int32(ccCVV)},
		UserId:       sessionID(r),
		UserCurrency: currentCurrency(r),
		Address: &pb.Address{
			StreetAddress: streetAddress,
			City:          city,
			State:         state,
			ZipCode:       int32(zipCode),
			Country:       country},
	})
	if err != nil {
		renderHTTPError(log, ctx, errors.Wrap(err, "下订单失败"), http.StatusInternalServerError)
		return
	}
	log.WithField("order", order.GetOrder().GetOrderId()).Info("下单")

	order.GetOrder().GetItems()
	recommendations, _ := fe.getRecommendations(r.Context(), sessionID(r), nil)

	totalPaid := order.GetOrder().GetShippingCost()
	for _, v := range order.GetOrder().GetItems() {
		multPrice := money.MultiplySlow(v.GetCost(), uint32(v.GetItem().GetQuantity()))
		totalPaid = money.Must(money.Sum(totalPaid, multPrice))
	}

	currencies, err := fe.getCurrencies(r.Context())
	if err != nil {
		renderHTTPError(log, ctx, errors.Wrap(err, "不能查询到货币"), http.StatusInternalServerError)
		return
	}

	resultMap := map[string]interface{}{
		"session_id":      sessionID(r),
		"request_id":      r.Context().Value(ctxKeyRequestID{}),
		"user_currency":   currentCurrency(r),
		"show_currency":   false,
		"currencies":      currencies,
		"order":           order.GetOrder(),
		"total_paid":      &totalPaid,
		"recommendations": recommendations,
	}

	ctx.HTML(http.StatusOK, "order", resultMap)
}

// 退出登录
func (fe *FrontendServer) logoutHandler(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	log.Debug("退出登录")
	for _, c := range r.Cookies() {
		c.Expires = time.Now().Add(-time.Hour * 24 * 365)
		c.MaxAge = -1
		http.SetCookie(w, c)
	}
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusFound)
}

// 设置货币
func (fe *FrontendServer) setCurrencyHandler(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer

	cur := r.FormValue("currency_code")
	log.WithField("curr.new", cur).WithField("curr.old", currentCurrency(r)).
		Debug("setting currency")

	if cur != "" {
		http.SetCookie(w, &http.Cookie{
			Name:   cookieCurrency,
			Value:  cur,
			MaxAge: cookieMaxAge,
		})
	}
	referer := r.Header.Get("referer")
	if referer == "" {
		referer = "/"
	}
	w.Header().Set("Location", referer)
	w.WriteHeader(http.StatusFound)
}

// 关闭广告
func (fe *FrontendServer) chooseAd(ctx context.Context, ctxKeys []string, log logrus.FieldLogger) *pb.Ad {
	ads, err := fe.getAd(ctx, ctxKeys)
	if err != nil {
		log.WithField("error", err).Warn("查询广告失败")
		return nil
	}
	if len(ads) == 0 {
		return nil
	}
	return ads[rand.Intn(len(ads))]
}

// 错误信息
func renderHTTPError(log logrus.FieldLogger, ctx *gin.Context, err error, code int) {
	r := ctx.Request
	w := ctx.Writer
	log.WithField("error", err).Error("请求错误")
	errMsg := fmt.Sprintf("%+v", err)

	w.WriteHeader(code)

	resultMap := map[string]interface{}{
		"session_id":  sessionID(r),
		"request_id":  r.Context().Value(ctxKeyRequestID{}),
		"error":       errMsg,
		"status_code": code,
		"status":      http.StatusText(code),
	}

	ctx.HTML(http.StatusOK, "error", resultMap)
}

// 当前货币
func currentCurrency(r *http.Request) string {
	c, _ := r.Cookie(cookieCurrency)
	if c != nil {
		return c.Value
	}
	return defaultCurrency
}

// session会话
func sessionID(r *http.Request) string {
	v := r.Context().Value(ctxKeySessionID{})
	if v != nil {
		return v.(string)
	}
	return ""
}

// 购物车id
func cartIDs(c []*pb.CartItem) []string {
	out := make([]string, len(c))
	for i, v := range c {
		out[i] = v.GetProductId()
	}
	return out
}

// cart size
func cartSize(c []*pb.CartItem) int {
	cartSize := 0
	for _, item := range c {
		cartSize += int(item.GetQuantity())
	}
	return cartSize
}

// 格式化货币
func renderMoney(money *pb.Money) string {
	currencyLogo := renderCurrencyLogo(money.GetCurrencyCode())
	return fmt.Sprintf("%s%d.%02d", currencyLogo, money.GetUnits(), money.GetNanos()/10000000)
}

// 货币符号
func renderCurrencyLogo(currencyCode string) string {
	logos := map[string]string{
		"USD": "$",
		"CAD": "$",
		"JPY": "¥",
		"EUR": "€",
		"TRY": "₺",
		"GBP": "£",
	}

	logo := "$" //default
	if val, ok := logos[currencyCode]; ok {
		logo = val
	}
	return logo
}

// 判断字符串是否在字符串切片
func stringinSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
