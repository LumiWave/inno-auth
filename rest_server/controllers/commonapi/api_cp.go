package commonapi

import (
	"net/http"

	"github.com/ONBUFF-IP-TOKEN/baseapp/base"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/resultcode"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/model"
	"github.com/labstack/echo"
)

func PostCPRegister(c echo.Context, company *context.Company) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// CP사이름 빈문자열 체크
	if err := company.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	// CP사이름 중복 체크
	if value, err := model.GetDB().SelectGetCpInfoByCpName(company.CompanyName); err == nil {
		if len(value.CompanyName) > 0 {
			log.Error("PostCPRegister exists cp=", company.CompanyName, " errorCode:", resultcode.Result_Auth_ExistsCpName)
			resp.SetReturn(resultcode.Result_Auth_ExistsCpName)
			return c.JSON(http.StatusOK, resp)
		}
	}

	// 테이블에 신규 row 생성
	if err := model.GetDB().InsertCP(company); err != nil {
		log.Error(err)
		resp.SetReturn(resultcode.Result_DBError)
	}

	return c.JSON(http.StatusOK, resp)
}

func DelCPUnRegister(c echo.Context, company *context.Company) error {
	resp := new(base.BaseResponse)
	resp.Success()

	// CP사이름 빈문자열 체크
	if err := company.CheckValidate(); err != nil {
		return c.JSON(http.StatusOK, err)
	}

	// 테이블 row 삭제
	if err := model.GetDB().DeleteCP(company); err != nil {
		log.Error(err)
		resp.SetReturn(resultcode.Result_DBError)
	}
	return c.JSON(http.StatusOK, resp)
}

func GetCPExists(c echo.Context, company *context.Company) error {
	resp := new(base.BaseResponse)
	resp.Success()

	if respCpInfo, err := model.GetDB().SelectGetCpInfoByCpName(company.CompanyName); err != nil {
		resp.SetReturn(resultcode.Result_DBError)
	} else {
		if len(respCpInfo.CompanyName) != 0 {
			resp.Value = respCpInfo
		} else {
			resp.SetReturn(resultcode.Result_Auth_NotFoundCpName)
		}
	}
	return c.JSON(http.StatusOK, resp)
}
