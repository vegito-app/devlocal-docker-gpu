package v1

type StoreUserRecoveryKeyRequestBody struct {
	UserID      string `json:"userId"`
	RecoveryKey []byte `json:"recoveryKey"`
}
