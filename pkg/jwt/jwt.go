package jwt

import (
	"errors"
	"time"

	"github.com/spf13/viper"

	"github.com/dgrijalva/jwt-go"
)

var mySecret = []byte("汤雄胜")

var ErrNeedLogin = errors.New("token过期,请重新登陆")

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	Email  string `json:"email"`
	UserID int64  `json:"user_id"`
	jwt.StandardClaims
}

// GenToken 生成JWT
func GenToken(userID int64, email string) (aToken, rToken string, err error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		Email:  email,
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().
				Add(time.Minute * 90).Unix(), // 过期时间
			Issuer: "manage", // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	//使用指定的secret签名并获得完整的编码后的字符串token
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(mySecret)

	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(viper.GetInt("auth.jwt_expire")) * time.Hour).Unix(),
			Issuer:    "manage",
		}).SignedString(mySecret)
	return
}

// ParseToken 解析JWT

func ParseToken(tokenString string) (*MyClaims, error) {
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid { // 校验token
		return mc, nil
	}
	return nil, errors.New("invalid token")
}

//RefreshToken 刷新AccessToken
func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	//refresh token 无效直接返回
	if _, err = jwt.Parse(rToken, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	}); err != nil {
		err = ErrNeedLogin
		return
	}
	//从旧access token 解析出claims数据
	var claims MyClaims
	_, err = jwt.ParseWithClaims(aToken, &claims, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	v, _ := err.(*jwt.ValidationError)
	// 当access token是过期错误 并且refresh token 没有过期就创建一个新的access token
	if v.Errors == jwt.ValidationErrorExpired {
		return GenToken(claims.UserID, claims.Email)
	}
	return
}
