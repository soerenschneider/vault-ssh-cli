package ssh

import "fmt"

type RefreshSignatureStrategy interface {
	NeedsNewSignature(*CertInfo) bool
}

type SimpleStrategy struct {
	renew bool
}

func NewSimpleStrategy(renew bool) *SimpleStrategy {
	return &SimpleStrategy{renew}
}

func (s *SimpleStrategy) NeedsNewSignature(lifetime *CertInfo) bool {
	return s.renew
}

type PercentageRenewStrategy struct {
	minPercentageLeft float32
}

func NewPercentageStrategy(minPercentageLeft float32) (*PercentageRenewStrategy, error) {
	if minPercentageLeft < 1 || minPercentageLeft > 99 {
		return nil, fmt.Errorf("invalid value: %f", minPercentageLeft)
	}

	return &PercentageRenewStrategy{minPercentageLeft: minPercentageLeft}, nil
}

func (s *PercentageRenewStrategy) NeedsNewSignature(lifetime *CertInfo) bool {
	return lifetime.GetPercentage() <= s.minPercentageLeft
}
