package internal

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	corev1 "k8s.io/api/core/v1"
	kubeErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type KubernetesSecret struct {
	client     *kubernetes.Clientset
	namespace  string
	secretName string
	secretKey  string
}

func NewKubernetesSecret(client *kubernetes.Clientset, secretName, secretKey, namespace string) (*KubernetesSecret, error) {
	if strings.TrimSpace(secretName) == "" {
		return nil, errors.New("empty secretName provided")
	}

	if strings.TrimSpace(secretKey) == "" {
		return nil, errors.New("empty secretKey provided")
	}

	if strings.TrimSpace(namespace) == "" {
		detectedNamespace, err := detectNamespace()
		if err != nil {
			return nil, fmt.Errorf("no namespace was given and could not auto-detect namespace: %w", err)
		}
		namespace = detectedNamespace
	}

	if client == nil {
		return nil, errors.New("empty kubernetes client passed")
	}

	return &KubernetesSecret{
		client:     client,
		namespace:  namespace,
		secretName: secretName,
		secretKey:  secretKey,
	}, nil
}

func detectNamespace() (string, error) {
	namespaceFile := "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
	namespace, err := os.ReadFile(namespaceFile)
	if err != nil {
		return "", fmt.Errorf("failed to read namespace from %s: %v", namespaceFile, err)
	}

	return string(namespace), nil
}

func (ks *KubernetesSecret) Read() ([]byte, error) {
	secret, err := ks.client.CoreV1().Secrets(ks.namespace).Get(context.Background(), ks.secretName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not get secret %s: %v", ks.secretName, err)
	}

	encodedValue, exists := secret.Data[ks.secretKey]
	if !exists {
		return nil, fmt.Errorf("secret key %s not found in secret %s", ks.secretKey, ks.secretName)
	}

	return encodedValue, nil
}

func (ks *KubernetesSecret) Write(value string) error {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: ks.secretName,
		},
		Data: map[string][]byte{
			"signature": []byte(value),
		},
	}

	_, err := ks.client.CoreV1().Secrets(ks.namespace).Get(context.TODO(), ks.secretName, metav1.GetOptions{})
	if err != nil {
		if kubeErrors.IsNotFound(err) {
			_, err = ks.client.CoreV1().Secrets(ks.namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
			return err
		} else if kubeErrors.IsForbidden(err) {
			return err
		}

		return err
	}

	_, err = ks.client.CoreV1().Secrets(ks.namespace).Update(context.Background(), secret, metav1.UpdateOptions{})
	return err
}
