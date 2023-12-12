package consumer

import (
	"fmt"
	"net/http"

	"github.com/free5gc/openapi/Nnrf_NFDiscovery"
	"github.com/free5gc/openapi/models"
	udm_context "github.com/free5gc/udm/internal/context"
	"github.com/free5gc/udm/internal/logger"
	"github.com/free5gc/udm/internal/util"
)

const (
	NFDiscoveryToUDRParamNone int = iota
	NFDiscoveryToUDRParamSupi
	NFDiscoveryToUDRParamExtGroupId
	NFDiscoveryToUDRParamGpsi
)

func SendNFIntances(nrfUri string, targetNfType, requestNfType models.NfType,
	param Nnrf_NFDiscovery.SearchNFInstancesParamOpts,
) (*models.SearchResult, error) {
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath(nrfUri) // addr
	clientNRF := Nnrf_NFDiscovery.NewAPIClient(configuration)

	ctx, _, err := GetTokenCtx("nnrf-disc", "NRF")
	if err != nil {
		return nil, err
	}

	result, res, err1 := clientNRF.NFInstancesStoreApi.SearchNFInstances(ctx, targetNfType,
		requestNfType, &param)
	if err1 != nil {
		err = err1
		return nil, err
	}

	defer func() {
		if rspCloseErr := res.Body.Close(); rspCloseErr != nil {
			logger.ConsumerLog.Errorf("SearchNFInstances response body cannot close: %+v", rspCloseErr)
		}
	}()
	if res != nil && res.StatusCode == http.StatusTemporaryRedirect {
		err = fmt.Errorf("Temporary Redirect For Non NRF Consumer")
	}
	return &result, err
}

func SendNFIntancesUDR(id string, types int) string {
	self := udm_context.GetSelf()
	targetNfType := models.NfType_UDR
	requestNfType := models.NfType_UDM
	localVarOptionals := Nnrf_NFDiscovery.SearchNFInstancesParamOpts{
		// 	DataSet: optional.NewInterface(models.DataSetId_SUBSCRIPTION),
	}
	// switch types {
	// case NFDiscoveryToUDRParamSupi:
	// 	localVarOptionals.Supi = optional.NewString(id)
	// case NFDiscoveryToUDRParamExtGroupId:
	// 	localVarOptionals.ExternalGroupIdentity = optional.NewString(id)
	// case NFDiscoveryToUDRParamGpsi:
	// 	localVarOptionals.Gpsi = optional.NewString(id)
	// }
	fmt.Println(self.NrfUri)
	result, err := SendNFIntances(self.NrfUri, targetNfType, requestNfType, localVarOptionals)
	if err != nil {
		logger.ConsumerLog.Error(err.Error())
		return ""
	}
	for _, profile := range result.NfInstances {
		return util.SearchNFServiceUri(profile, models.ServiceName_NUDR_DR, models.NfServiceStatus_REGISTERED)
	}
	return ""
}
