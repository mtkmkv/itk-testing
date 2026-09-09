package domain

type OperationType string

const (
	OperationTypeDeposit  OperationType = "DEPOSIT"
	OperationTypeWithdraw OperationType = "WITHDRAW"
)

func (o OperationType) IsValid() bool {
	switch o {
	case OperationTypeDeposit, OperationTypeWithdraw:
		return true
	default:
		return false
	}
}