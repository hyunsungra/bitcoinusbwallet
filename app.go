package main

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	goruntime "runtime"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/tyler-smith/go-bip39"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/text/unicode/norm"
)

// App struct
type App struct {
	ctx context.Context
}

// WalletData 지갑 정보를 저장하는 구조체 (coldwallet 호환)
type WalletData struct {
	Name          string `json:"name"`          // 지갑 이름
	Mnemonic      string `json:"mnemonic"`      // 니모닉 구문
	Passphrase    string `json:"passphrase"`    // 추가 패스프레이즈
	Address       string `json:"address"`       // 비트코인 주소
	PublicKey     string `json:"publicKey"`     // 공개키
	PrivateKeyWIF string `json:"privateKeyWIF"` // WIF 형식 개인키
	Path          string `json:"path"`          // BIP 파생 경로
	CreatedAt     string `json:"createdAt"`     // 생성 시간
}

// ColdWalletFileFormat coldwallet 파일 형식
type ColdWalletFileFormat struct {
	Version       string                 `json:"version"`
	Algorithm     string                 `json:"algorithm"`
	KeyDerivation string                 `json:"keyDerivation"`
	PBKDF2Params  map[string]interface{} `json:"pbkdf2Params"`
	Salt          string                 `json:"salt"`
	IV            string                 `json:"iv"`
	Data          string                 `json:"data"`
	Checksum      string                 `json:"checksum"`
}

// CreateWalletRequest 지갑 생성 요청 구조체
type CreateWalletRequest struct {
	Name       string `json:"name"`       // 지갑 이름
	Password   string `json:"password"`   // 지갑 비밀번호
	Mnemonic   string `json:"mnemonic"`   // 니모닉 구문
	Passphrase string `json:"passphrase"` // 추가 패스프레이즈
	SavePath   string `json:"savePath"`   // 저장 경로
}

// CheckWalletRequest 지갑 확인 요청 구조체
type CheckWalletRequest struct {
	FilePath string `json:"filePath"` // 지갑 파일 경로
	Password string `json:"password"` // 지갑 비밀번호
}

// CreateWalletResponse 지갑 생성 응답 구조체
type CreateWalletResponse struct {
	Success  bool   `json:"success"`            // 성공 여부
	Message  string `json:"message"`            // 응답 메시지
	FilePath string `json:"filePath,omitempty"` // 저장된 파일 경로
}

// PasswordValidation 비밀번호 검증 결과 구조체
type PasswordValidation struct {
	IsValid  bool     `json:"isValid"`  // 유효성 검사 통과 여부
	Errors   []string `json:"errors"`   // 오류 메시지 목록
	Strength string   `json:"strength"` // 비밀번호 강도
	Score    int      `json:"score"`    // 점수 (0-7)
}

// NewApp 새로운 App 애플리케이션 구조체 생성
func NewApp() *App {
	return &App{}
}

// Startup 앱 시작시 호출되는 함수, 컨텍스트 저장
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet 인사말 반환 함수 (테스트용)
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GenerateMnemonic BIP39 표준을 사용하여 24단어 니모닉 구문 생성
func (a *App) GenerateMnemonic() string {
	entropy, err := bip39.NewEntropy(256) // 256비트 엔트로피로 24단어 생성
	if err != nil {
		return ""
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return ""
	}

	return mnemonic
}

// ValidatePassword 비밀번호 강도 및 보안 요구사항 검증
func (a *App) ValidatePassword(password string) PasswordValidation {
	validation := PasswordValidation{
		IsValid: true,
		Errors:  []string{},
		Score:   0,
	}

	// 최소 길이 검사 (8자 이상)
	if len(password) < 8 {
		validation.IsValid = false
		validation.Errors = append(validation.Errors, "비밀번호는 8자 이상이어야 합니다")
	} else {
		validation.Score += 1
	}

	// 대문자 포함 검사
	hasUpper := false
	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
			break
		}
	}
	if !hasUpper {
		validation.IsValid = false
		validation.Errors = append(validation.Errors, "대문자를 최소 1개 포함해야 합니다")
	} else {
		validation.Score += 1
	}

	// 소문자 포함 검사
	hasLower := false
	for _, char := range password {
		if unicode.IsLower(char) {
			hasLower = true
			break
		}
	}
	if !hasLower {
		validation.IsValid = false
		validation.Errors = append(validation.Errors, "소문자를 최소 1개 포함해야 합니다")
	} else {
		validation.Score += 1
	}

	// 숫자 포함 검사
	hasDigit := false
	for _, char := range password {
		if unicode.IsDigit(char) {
			hasDigit = true
			break
		}
	}
	if !hasDigit {
		validation.IsValid = false
		validation.Errors = append(validation.Errors, "숫자를 최소 1개 포함해야 합니다")
	} else {
		validation.Score += 1
	}

	// 특수문자 포함 검사
	hasSpecial := false
	for _, char := range password {
		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSpecial = true
			break
		}
	}
	if !hasSpecial {
		validation.IsValid = false
		validation.Errors = append(validation.Errors, "특수문자를 최소 1개 포함해야 합니다")
	} else {
		validation.Score += 1
	}

	// 연속된 문자 검사 (abc, 123 등)
	if hasConsecutiveChars(password) {
		validation.IsValid = false
		validation.Errors = append(validation.Errors, "연속된 문자는 사용할 수 없습니다")
	} else {
		validation.Score += 1
	}

	// 반복된 문자 검사 (aaa, 111 등)
	if hasRepeatedChars(password) {
		validation.IsValid = false
		validation.Errors = append(validation.Errors, "같은 문자를 3번 이상 연속으로 사용할 수 없습니다")
	} else {
		validation.Score += 1
	}

	// 비밀번호 강도 계산
	if validation.Score >= 6 {
		validation.Strength = "매우 강함"
	} else if validation.Score >= 5 {
		validation.Strength = "강함"
	} else if validation.Score >= 4 {
		validation.Strength = "보통"
	} else if validation.Score >= 2 {
		validation.Strength = "약함"
	} else {
		validation.Strength = "매우 약함"
	}

	return validation
}

// hasConsecutiveChars 비밀번호에 연속된 문자가 있는지 검사
func hasConsecutiveChars(password string) bool {
	password = strings.ToLower(password)
	for i := 0; i < len(password)-2; i++ {
		if password[i]+1 == password[i+1] && password[i+1]+1 == password[i+2] {
			return true
		}
	}
	return false
}

// hasRepeatedChars 비밀번호에 반복된 문자가 있는지 검사
func hasRepeatedChars(password string) bool {
	for i := 0; i < len(password)-2; i++ {
		if password[i] == password[i+1] && password[i+1] == password[i+2] {
			return true
		}
	}
	return false
}

// CreateWallet 새로운 비트코인 지갑 생성 및 저장
func (a *App) CreateWallet(request CreateWalletRequest) CreateWalletResponse {
	// 니모닉 유효성 검증 제거 (부분 수정 허용)
	// coldwallet도 체크섬 검증을 우회할 수 있음
	/*
		if !bip39.IsMnemonicValid(request.Mnemonic) {
			return CreateWalletResponse{
				Success: false,
				Message: "유효하지 않은 니모닉 구문입니다",
			}
		}
	*/

	// 패스프레이즈를 NFKD 정규화 (한글 패스프레이즈 호환)
	passphraseNormalized := string(norm.NFKD.Bytes([]byte(request.Passphrase)))

	// 니모닉과 정규화된 패스프레이즈로부터 시드 생성
	seed := bip39.NewSeed(request.Mnemonic, passphraseNormalized)

	// 마스터 키 생성
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return CreateWalletResponse{
			Success: false,
			Message: "마스터 키 생성 실패: " + err.Error(),
		}
	}

	// BIP84 경로를 사용한 키 파생: m/84'/0'/0'/0/0 (coldwallet과 동일)
	purpose, err := masterKey.Derive(hdkeychain.HardenedKeyStart + 84)
	if err != nil {
		return CreateWalletResponse{
			Success: false,
			Message: "BIP84 purpose 파생 실패: " + err.Error(),
		}
	}

	coinType, err := purpose.Derive(hdkeychain.HardenedKeyStart + 0)
	if err != nil {
		return CreateWalletResponse{
			Success: false,
			Message: "Coin type 파생 실패: " + err.Error(),
		}
	}

	account, err := coinType.Derive(hdkeychain.HardenedKeyStart + 0)
	if err != nil {
		return CreateWalletResponse{
			Success: false,
			Message: "Account 파생 실패: " + err.Error(),
		}
	}

	change, err := account.Derive(0)
	if err != nil {
		return CreateWalletResponse{
			Success: false,
			Message: "Change 파생 실패: " + err.Error(),
		}
	}

	addressKey, err := change.Derive(0)
	if err != nil {
		return CreateWalletResponse{
			Success: false,
			Message: "Address key 파생 실패: " + err.Error(),
		}
	}

	// 개인키 추출
	privateKey, err := addressKey.ECPrivKey()
	if err != nil {
		return CreateWalletResponse{
			Success: false,
			Message: "개인키 생성 실패: " + err.Error(),
		}
	}

	// 공개키 추출
	publicKey := privateKey.PubKey()

	// SegWit 주소를 위한 witness program 생성 (coldwallet과 동일)
	witnessProgram := btcutil.Hash160(publicKey.SerializeCompressed())

	// Native SegWit 주소 생성 (coldwallet: bitcoin.payments.p2wpkh)
	address, err := btcutil.NewAddressWitnessPubKeyHash(witnessProgram, &chaincfg.MainNetParams)
	if err != nil {
		return CreateWalletResponse{
			Success: false,
			Message: "주소 생성 실패: " + err.Error(),
		}
	}

	// WIF 개인키 생성
	privateKeyWIF, err := btcutil.NewWIF(privateKey, &chaincfg.MainNetParams, true)
	if err != nil {
		return CreateWalletResponse{
			Success: false,
			Message: "WIF 개인키 생성 실패: " + err.Error(),
		}
	}

	// 지갑 데이터 구조체 생성 (coldwallet 호환 형식)
	walletData := WalletData{
		Name:          request.Name,
		Mnemonic:      request.Mnemonic,
		Passphrase:    request.Passphrase,
		Address:       address.EncodeAddress(),
		PublicKey:     hex.EncodeToString(publicKey.SerializeCompressed()),
		PrivateKeyWIF: privateKeyWIF.String(),
		Path:          "m/84'/0'/0'/0/0",
		CreatedAt:     time.Now().Format(time.RFC3339),
	}

	// 지갑 암호화 및 저장 (coldwallet 호환 방식)
	filePath, err := a.saveColdWallet(walletData, request.Password, request.Name, request.SavePath)
	if err != nil {
		return CreateWalletResponse{
			Success: false,
			Message: "지갑 저장 실패: " + err.Error(),
		}
	}

	return CreateWalletResponse{
		Success:  true,
		Message:  "지갑이 성공적으로 생성되었습니다",
		FilePath: filePath,
	}
}

// derivePath 파생 경로 문자열을 파싱하여 키 파생
func derivePath(masterKey *hdkeychain.ExtendedKey, path string) (*hdkeychain.ExtendedKey, error) {
	if !strings.HasPrefix(path, "m/") {
		return nil, fmt.Errorf("invalid path: %s", path)
	}

	// "m/" 제거
	path = strings.TrimPrefix(path, "m/")
	if path == "" {
		return masterKey, nil
	}

	// 각 세그먼트 파싱
	segments := strings.Split(path, "/")
	currentKey := masterKey

	for _, segment := range segments {
		if segment == "" {
			continue
		}

		var index uint32
		hardened := false

		if strings.HasSuffix(segment, "'") || strings.HasSuffix(segment, "h") {
			hardened = true
			segment = strings.TrimSuffix(strings.TrimSuffix(segment, "'"), "h")
		}

		// 문자열을 숫자로 변환
		val, err := fmt.Sscanf(segment, "%d", &index)
		if err != nil || val != 1 {
			return nil, fmt.Errorf("invalid path segment: %s", segment)
		}

		if hardened {
			index += hdkeychain.HardenedKeyStart
		}

		currentKey, err = currentKey.Derive(index)
		if err != nil {
			return nil, fmt.Errorf("failed to derive key at index %d: %v", index, err)
		}
	}

	return currentKey, nil
}

// saveColdWallet coldwallet 호환 방식으로 지갑 데이터 암호화 및 저장
func (a *App) saveColdWallet(walletData WalletData, password, walletName, savePath string) (string, error) {
	// 지갑 데이터를 JSON으로 직렬화
	walletJSON, err := json.Marshal(walletData)
	if err != nil {
		return "", err
	}

	// 체크섬 계산
	checksum := sha256.Sum256(walletJSON)

	// Salt 및 IV 생성 (32바이트 salt, 16바이트 IV)
	salt := make([]byte, 32)
	iv := make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		return "", err
	}
	_, err = rand.Read(iv)
	if err != nil {
		return "", err
	}

	// PBKDF2로 키 생성 (coldwallet와 동일)
	key := pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)

	// Node.js crypto.createCipher('aes-256-cbc', key) 에뮬레이션
	// deprecated createCipher는 EVP_BytesToKey 방식으로 키를 파생하는데
	// 여기서는 이미 PBKDF2로 파생된 키를 사용

	// AES 암호화를 위한 cipher 생성
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// PKCS7 패딩 추가
	paddedData := pkcs7Pad(walletJSON, aes.BlockSize)

	// createCipher는 내부적으로 랜덤 IV를 생성하지만
	// 여기서는 미리 생성된 IV를 사용
	mode := cipher.NewCBCEncrypter(block, iv)

	// CBC 모드로 암호화
	encrypted := make([]byte, len(paddedData))
	mode.CryptBlocks(encrypted, paddedData)

	// hex 인코딩 (데이터 필드에 저장될 값)
	encryptedHex := hex.EncodeToString(encrypted)

	// coldwallet 파일 형식 생성
	fileFormat := ColdWalletFileFormat{
		Version:       "2.0",
		Algorithm:     "aes-256-cbc",
		KeyDerivation: "pbkdf2",
		PBKDF2Params: map[string]interface{}{
			"iterations": 100000,
			"digest":     "sha256",
		},
		Salt:     hex.EncodeToString(salt),
		IV:       hex.EncodeToString(iv),
		Data:     encryptedHex,
		Checksum: hex.EncodeToString(checksum[:]),
	}

	// JSON으로 직렬화
	fileData, err := json.Marshal(fileFormat)
	if err != nil {
		return "", err
	}

	// 저장 경로 결정
	var saveDir string
	if savePath != "" {
		saveDir = savePath
	} else {
		// 기본 저장 경로 (폴더를 선택하지 않은 경우)
		if goruntime.GOOS == "windows" {
			saveDir = "C:\\"
		} else {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				saveDir = "."
			} else {
				saveDir = homeDir
			}
		}
	}

	// 파일명 생성 (공백을 언더스코어로 치환 및 중복 방지)
	baseName := strings.ReplaceAll(walletName, " ", "_")
	filename := baseName + ".wallet"
	filePath := filepath.Join(saveDir, filename)

	// 파일이 이미 존재하면 번호를 추가하여 중복 방지
	counter := 1
	for {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// 파일이 존재하지 않으면 사용 가능
			break
		}
		// 파일이 존재하면 번호 추가
		filename = fmt.Sprintf("%s(%d).wallet", baseName, counter)
		filePath = filepath.Join(saveDir, filename)
		counter++
	}

	// 파일 저장 (소유자만 읽기/쓰기 권한)
	err = os.WriteFile(filePath, fileData, 0600)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

// pkcs7Pad 데이터에 PKCS7 패딩을 추가
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	return append(data, padtext...)
}

// GetBIP39WordList BIP39 표준 단어 목록 반환 (2048개)
func (a *App) GetBIP39WordList() []string {
	return bip39.GetWordList()
}

// CheckWalletResponse 지갑 확인 응답 구조체
type CheckWalletResponse struct {
	Success    bool       `json:"success"`    // 성공 여부
	Message    string     `json:"message"`    // 오류 메시지
	WalletData WalletData `json:"walletData"` // 지갑 데이터
}

// OpenWalletResponse 지갑 열기 응답 구조체 (비트코인 전송용)
type OpenWalletResponse struct {
	Success    bool   `json:"success"`    // 성공 여부
	Message    string `json:"message"`    // 응답 메시지
	Address    string `json:"address"`    // 비트코인 주소
	PrivateKey string `json:"privateKey"` // 개인키 (WIF 형식)
}

// CheckWallet 지갑 파일 검증 및 데이터 반환
func (a *App) CheckWallet(request CheckWalletRequest) CheckWalletResponse {
	// 파일 존재 확인
	if _, err := os.Stat(request.FilePath); os.IsNotExist(err) {
		return CheckWalletResponse{
			Success: false,
			Message: "지갑 파일을 찾을 수 없습니다.",
		}
	}

	// 파일 읽기
	encryptedData, err := os.ReadFile(request.FilePath)
	if err != nil {
		return CheckWalletResponse{
			Success: false,
			Message: "지갑 파일을 읽을 수 없습니다: " + err.Error(),
		}
	}

	// 비밀번호로 복호화 (coldwallet 호환 방식)
	walletData, err := a.decryptColdWallet(encryptedData, request.Password)
	if err != nil {
		return CheckWalletResponse{
			Success: false,
			Message: "잘못된 비밀번호이거나 손상된 지갑 파일입니다.",
		}
	}

	return CheckWalletResponse{
		Success:    true,
		Message:    "성공",
		WalletData: walletData,
	}
}

// OpenWallet 지갑 파일을 열어서 주소와 개인키만 반환 (비트코인 전송용)
func (a *App) OpenWallet(request CheckWalletRequest) OpenWalletResponse {
	// 파일 존재 확인
	if _, err := os.Stat(request.FilePath); os.IsNotExist(err) {
		return OpenWalletResponse{
			Success: false,
			Message: "지갑 파일을 찾을 수 없습니다.",
		}
	}

	// 파일 읽기
	fileData, err := os.ReadFile(request.FilePath)
	if err != nil {
		return OpenWalletResponse{
			Success: false,
			Message: "지갑 파일을 읽을 수 없습니다: " + err.Error(),
		}
	}

	// 지갑 데이터 복호화
	walletData, err := a.decryptColdWallet(fileData, request.Password)
	if err != nil {
		return OpenWalletResponse{
			Success: false,
			Message: "잘못된 비밀번호이거나 손상된 지갑 파일입니다.",
		}
	}

	// 필요한 정보만 반환 (주소와 개인키)
	return OpenWalletResponse{
		Success:    true,
		Message:    "성공",
		Address:    walletData.Address,
		PrivateKey: walletData.PrivateKeyWIF,
	}
}

// decryptColdWallet coldwallet 호환 방식으로 지갑 데이터 복호화
func (a *App) decryptColdWallet(fileData []byte, password string) (WalletData, error) {
	// coldwallet 파일 형식 파싱
	var fileFormat ColdWalletFileFormat
	err := json.Unmarshal(fileData, &fileFormat)
	if err != nil {
		return WalletData{}, fmt.Errorf("잘못된 파일 형식: %v", err)
	}

	// 버전 확인
	if fileFormat.Version != "2.0" {
		return WalletData{}, fmt.Errorf("지원되지 않는 지갑 버전: %s", fileFormat.Version)
	}

	// Salt 및 암호화된 데이터 디코딩
	salt, err := hex.DecodeString(fileFormat.Salt)
	if err != nil {
		return WalletData{}, fmt.Errorf("Salt 디코딩 실패: %v", err)
	}

	encrypted, err := hex.DecodeString(fileFormat.Data)
	if err != nil {
		return WalletData{}, fmt.Errorf("암호화 데이터 디코딩 실패: %v", err)
	}

	// PBKDF2로 키 생성
	key := pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)

	// Node.js createDecipher('aes-256-cbc', key) 에뮬레이션
	// createDecipher는 키를 직접 사용
	block, err := aes.NewCipher(key)
	if err != nil {
		return WalletData{}, err
	}

	// 블록 크기 확인
	if len(encrypted)%aes.BlockSize != 0 {
		return WalletData{}, fmt.Errorf("잘못된 비밀번호이거나 손상된 파일입니다")
	}

	// 파일에 저장된 IV 사용
	iv, err := hex.DecodeString(fileFormat.IV)
	if err != nil {
		return WalletData{}, fmt.Errorf("IV 디코딩 실패: %v", err)
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	plaintext := make([]byte, len(encrypted))
	mode.CryptBlocks(plaintext, encrypted)

	// PKCS7 패딩 제거
	plaintext, err = pkcs7Unpad(plaintext)
	if err != nil {
		return WalletData{}, fmt.Errorf("잘못된 비밀번호입니다")
	}

	// JSON 파싱
	var walletData WalletData
	err = json.Unmarshal(plaintext, &walletData)
	if err != nil {
		return WalletData{}, fmt.Errorf("지갑 데이터 파싱 실패: %v", err)
	}

	// 체크섬 검증 (선택사항)
	if fileFormat.Checksum != "" {
		expectedChecksum, err := hex.DecodeString(fileFormat.Checksum)
		if err == nil {
			actualChecksum := sha256.Sum256(plaintext)
			for i, b := range actualChecksum {
				if i >= len(expectedChecksum) || b != expectedChecksum[i] {
					return WalletData{}, fmt.Errorf("지갑 데이터 무결성 검증 실패")
				}
			}
		}
	}

	return walletData, nil
}

// pkcs7Unpad PKCS7 패딩 제거
func pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("빈 데이터")
	}

	padding := int(data[length-1])
	if padding > length {
		return nil, fmt.Errorf("잘못된 패딩")
	}

	return data[:length-padding], nil
}

// SelectSaveDirectory 지갑 파일을 저장할 디렉터리 선택 대화상자 표시
func (a *App) SelectSaveDirectory() (string, error) {
	// 운영체제별 기본 경로 설정
	var defaultPath string
	if goruntime.GOOS == "windows" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			defaultPath = "C:\\"
		} else {
			defaultPath = filepath.Join(homeDir, "Desktop")
		}
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			defaultPath = "."
		} else {
			defaultPath = homeDir
		}
	}

	// Wails 런타임을 사용한 폴더 선택 다이얼로그
	selectedPath, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            "지갑 파일을 저장할 폴더 선택",
		DefaultDirectory: defaultPath,
	})

	if err != nil {
		// 오류 발생 시 기본 경로 반환
		return defaultPath, err
	}

	// 사용자가 취소를 선택한 경우 (빈 문자열 반환)
	if selectedPath == "" {
		return "", fmt.Errorf("폴더 선택이 취소되었습니다")
	}

	return selectedPath, nil
}

// SelectWalletFile 지갑 파일 선택 대화상자 표시
func (a *App) SelectWalletFile() (string, error) {
	// 지갑 파일 선택 대화상자
	selectedPath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "지갑 파일 선택",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "지갑 파일 (*.wallet)",
				Pattern:     "*.wallet",
			},
			{
				DisplayName: "모든 파일 (*.*)",
				Pattern:     "*.*",
			},
		},
	})

	if err != nil {
		return "", err
	}

	// 사용자가 취소를 선택한 경우
	if selectedPath == "" {
		return "", fmt.Errorf("파일 선택이 취소되었습니다")
	}

	return selectedPath, nil
}

// GetBalanceRequest 잔액 조회 요청 구조체
type GetBalanceRequest struct {
	Address string `json:"address"` // 비트코인 주소
}

// GetBalanceResponse 잔액 조회 응답 구조체
type GetBalanceResponse struct {
	Success     bool    `json:"success"`     // 성공 여부
	Message     string  `json:"message"`     // 응답 메시지
	Balance     float64 `json:"balance"`     // 잔액 (BTC)
	BalanceSat  int64   `json:"balanceSat"`  // 잔액 (satoshi)
	Confirmed   float64 `json:"confirmed"`   // 확인된 잔액 (BTC)
	Unconfirmed float64 `json:"unconfirmed"` // 미확인 잔액 (BTC)
	UTXOCount   int     `json:"utxoCount"`   // UTXO 개수
}

// AddressStats 주소 통계 정보 (Blockstream API)
type AddressStats struct {
	ChainStats struct {
		FundedTxoCount int64 `json:"funded_txo_count"`
		FundedTxoSum   int64 `json:"funded_txo_sum"`
		SpentTxoCount  int64 `json:"spent_txo_count"`
		SpentTxoSum    int64 `json:"spent_txo_sum"`
		TxCount        int64 `json:"tx_count"`
	} `json:"chain_stats"`
	MempoolStats struct {
		FundedTxoCount int64 `json:"funded_txo_count"`
		FundedTxoSum   int64 `json:"funded_txo_sum"`
		SpentTxoCount  int64 `json:"spent_txo_count"`
		SpentTxoSum    int64 `json:"spent_txo_sum"`
		TxCount        int64 `json:"tx_count"`
	} `json:"mempool_stats"`
}

// GetBalance 주소의 잔액 조회
func (a *App) GetBalance(request GetBalanceRequest) GetBalanceResponse {
	if request.Address == "" {
		return GetBalanceResponse{
			Success: false,
			Message: "주소를 입력해주세요",
		}
	}

	// 1. 주소 통계 조회
	statsUrl := fmt.Sprintf("https://blockstream.info/api/address/%s", request.Address)

	resp, err := http.Get(statsUrl)
	if err != nil {
		return GetBalanceResponse{
			Success: false,
			Message: fmt.Sprintf("주소 통계 조회 실패: %v", err),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return GetBalanceResponse{
			Success: false,
			Message: fmt.Sprintf("주소 통계 API 오류: %d", resp.StatusCode),
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return GetBalanceResponse{
			Success: false,
			Message: fmt.Sprintf("응답 읽기 실패: %v", err),
		}
	}

	var stats AddressStats
	if err := json.Unmarshal(body, &stats); err != nil {
		return GetBalanceResponse{
			Success: false,
			Message: fmt.Sprintf("주소 통계 파싱 실패: %v", err),
		}
	}

	// 2. UTXO 조회하여 확인된 잔액과 미확인 잔액 계산
	utxos, err := a.fetchUTXOs(request.Address)
	if err != nil {
		return GetBalanceResponse{
			Success: false,
			Message: fmt.Sprintf("UTXO 조회 실패: %v", err),
		}
	}

	// 확인된 잔액 계산
	var confirmedBalance int64
	var unconfirmedBalance int64

	// 모든 UTXO 조회 (확인된 것과 미확인된 것 모두)
	allUtxosUrl := fmt.Sprintf("https://blockstream.info/api/address/%s/utxo", request.Address)
	resp, err = http.Get(allUtxosUrl)
	if err == nil && resp.StatusCode == 200 {
		body, err := io.ReadAll(resp.Body)
		if err == nil {
			var allUtxos []UTXO
			if json.Unmarshal(body, &allUtxos) == nil {
				for _, utxo := range allUtxos {
					if utxo.Status.Confirmed {
						confirmedBalance += utxo.Value
					} else {
						unconfirmedBalance += utxo.Value
					}
				}
			}
		}
		resp.Body.Close()
	}

	// 총 잔액 계산
	totalBalance := confirmedBalance + unconfirmedBalance

	return GetBalanceResponse{
		Success:     true,
		Message:     "잔액 조회 성공",
		Balance:     float64(totalBalance) / 100000000,       // BTC 단위
		BalanceSat:  totalBalance,                            // satoshi 단위
		Confirmed:   float64(confirmedBalance) / 100000000,   // 확인된 잔액 (BTC)
		Unconfirmed: float64(unconfirmedBalance) / 100000000, // 미확인 잔액 (BTC)
		UTXOCount:   len(utxos),                              // 확인된 UTXO 개수
	}
}

// SendBitcoinRequest 비트코인 전송 요청 구조체
type SendBitcoinRequest struct {
	WalletData                WalletData `json:"walletData"`                // 지갑 데이터
	RecipientAddress          string     `json:"recipientAddress"`          // 받는 주소
	Amount                    float64    `json:"amount"`                    // 전송 금액 (BTC)
	FeeSatoshi                int        `json:"feeSatoshi"`                // 수수료 (사토시)
	IsDeveloperFeeTransaction bool       `json:"isDeveloperFeeTransaction"` // 개발자 수수료 트랜잭션 여부
	EnableFeeSplit            bool       `json:"enableFeeSplit"`            // 수수료 분할 활성화 여부
	DeveloperAddress          string     `json:"developerAddress"`          // 개발자 비트코인 주소
	DeveloperFeeSatoshi       int        `json:"developerFeeSatoshi"`       // 개발자 수수료 (사토시)
}

// SendBitcoinResponse 비트코인 전송 응답 구조체
type SendBitcoinResponse struct {
	Success   bool   `json:"success"`   // 성공 여부
	Message   string `json:"message"`   // 응답 메시지
	ErrorCode string `json:"errorCode"` // 에러 코드 (다국어 처리용)
	TxHash    string `json:"txHash"`    // 거래 해시
}

// UTXO 비트코인 UTXO 정보 구조체
type UTXO struct {
	TxID   string `json:"txid"`
	Vout   int    `json:"vout"`
	Value  int64  `json:"value"`
	Status struct {
		Confirmed bool `json:"confirmed"`
	} `json:"status"`
}

// TxOutput 거래 출력 정보
type TxOutput struct {
	ScriptPubKey string `json:"scriptpubkey"`
	Value        int64  `json:"value"`
}

// TxDetails 거래 세부정보
type TxDetails struct {
	Vout []TxOutput `json:"vout"`
}

// fetchUTXOs 주소의 UTXO 조회 (Blockstream API 사용)
func (a *App) fetchUTXOs(address string) ([]UTXO, error) {
	url := fmt.Sprintf("https://blockstream.info/api/address/%s/utxo", address)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("UTXO 조회 실패: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("UTXO 조회 API 오류: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("응답 읽기 실패: %v", err)
	}

	var utxos []UTXO
	if err := json.Unmarshal(body, &utxos); err != nil {
		return nil, fmt.Errorf("UTXO 파싱 실패: %v", err)
	}

	// 확인된 UTXO만 필터링
	var confirmedUTXOs []UTXO
	for _, utxo := range utxos {
		if utxo.Status.Confirmed {
			confirmedUTXOs = append(confirmedUTXOs, utxo)
		}
	}

	return confirmedUTXOs, nil
}

// fetchTxDetails 거래 세부정보 조회
func (a *App) fetchTxDetails(txid string) (*TxDetails, error) {
	url := fmt.Sprintf("https://blockstream.info/api/tx/%s", txid)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("거래 세부정보 조회 실패: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("거래 세부정보 API 오류: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("응답 읽기 실패: %v", err)
	}

	var txDetails TxDetails
	if err := json.Unmarshal(body, &txDetails); err != nil {
		return nil, fmt.Errorf("거래 세부정보 파싱 실패: %v", err)
	}

	return &txDetails, nil
}

// broadcastTransaction 거래 브로드캐스트
func (a *App) broadcastTransaction(txHex string) (string, error) {
	url := "https://blockstream.info/api/tx"

	resp, err := http.Post(url, "text/plain", strings.NewReader(txHex))
	if err != nil {
		return "", fmt.Errorf("거래 브로드캐스트 실패: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("응답 읽기 실패: %v", err)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("브로드캐스트 오류: %s", string(body))
	}

	return string(body), nil
}

// SendBitcoinTransaction 실제 비트코인 거래 전송
func (a *App) SendBitcoinTransaction(request SendBitcoinRequest) SendBitcoinResponse {
	// 요청 데이터 로깅
	/*
		fmt.Printf("=== SendBitcoinTransaction 호출 ===\n")
		fmt.Printf("받는 주소: %s\n", request.RecipientAddress)
		fmt.Printf("전송 금액: %.8f BTC\n", request.Amount)
		fmt.Printf("수수료: %d satoshi\n", request.FeeSatoshi)
		fmt.Printf("수수료 분할 활성화: %t\n", request.EnableFeeSplit)
		if request.EnableFeeSplit {
			fmt.Printf("개발자 주소: %s\n", request.DeveloperAddress)
			fmt.Printf("개발자 수수료: %d satoshi\n", request.DeveloperFeeSatoshi)
		}
		fmt.Printf("지갑 주소: %s\n", request.WalletData.Address)
		fmt.Printf("=====================================\n")
	*/

	// 입력 값 검증
	if request.RecipientAddress == "" {
		return SendBitcoinResponse{
			Success: false,
			Message: "받는 주소를 입력해주세요",
		}
	}

	if request.Amount <= 0 {
		return SendBitcoinResponse{
			Success: false,
			Message: "전송 금액은 0보다 커야 합니다",
		}
	}

	// 전체 수수료 범위 검증 (최소/최대) - 항상 체크
	if int64(request.FeeSatoshi) < 2000 {
		return SendBitcoinResponse{
			Success:   false,
			Message:   "전체 수수료가 너무 낮습니다. 최소 2000 사토시가 필요합니다.",
			ErrorCode: "FEE_TOO_LOW",
		}
	}
	if int64(request.FeeSatoshi) > 50000 {
		return SendBitcoinResponse{
			Success:   false,
			Message:   "전체 수수료가 너무 높습니다. 최대 50000 사토시를 초과할 수 없습니다.",
			ErrorCode: "FEE_TOO_HIGH",
		}
	}

	// UTXO 조회 전 수수료 분할 시스템 검증
	if request.EnableFeeSplit {

		// 채굴자 수수료 범위 검증 (최소/최대)
		minerFeeCheck := int64(request.FeeSatoshi) - int64(request.DeveloperFeeSatoshi)
		if minerFeeCheck < 1000 {
			return SendBitcoinResponse{
				Success:   false,
				Message:   "채굴자 수수료가 너무 낮습니다. 최소 1000 사토시가 필요합니다.",
				ErrorCode: "MINER_FEE_TOO_LOW",
			}
		}
		if minerFeeCheck > 50000 {
			return SendBitcoinResponse{
				Success:   false,
				Message:   "채굴자 수수료가 너무 높습니다. 최대 50000 사토시를 초과할 수 없습니다.",
				ErrorCode: "MINER_FEE_TOO_HIGH",
			}
		}

		// 개발자 수수료 범위 검증 (최소/최대)
		if request.DeveloperFeeSatoshi <= 0 {
			return SendBitcoinResponse{
				Success:   false,
				Message:   "개발자 수수료는 0보다 커야 합니다.",
				ErrorCode: "DEVELOPER_FEE_INVALID",
			}
		}
		if request.DeveloperFeeSatoshi > 10000 {
			return SendBitcoinResponse{
				Success:   false,
				Message:   "개발자 수수료가 너무 높습니다. 최대 10000 사토시를 초과할 수 없습니다.",
				ErrorCode: "DEVELOPER_FEE_TOO_HIGH",
			}
		}

		// 개발자 주소 검증
		if request.DeveloperAddress == "" {
			return SendBitcoinResponse{
				Success:   false,
				Message:   "개발자 주소를 입력해주세요.",
				ErrorCode: "DEVELOPER_ADDRESS_EMPTY",
			}
		}
	}

	// 더스트 한도 검증 (546 사토시)
	amountSatoshiCheck := int64(request.Amount * 100000000)
	if amountSatoshiCheck < 546 {
		return SendBitcoinResponse{
			Success:   false,
			Message:   "전송 금액이 너무 작습니다. 최소 546 사토시 (0.00000546 BTC)가 필요합니다.",
			ErrorCode: "AMOUNT_TOO_SMALL",
		}
	}

	// 1. UTXO 조회
	utxos, err := a.fetchUTXOs(request.WalletData.Address)
	if err != nil {
		return SendBitcoinResponse{
			Success: false,
			Message: fmt.Sprintf("UTXO 조회 실패: %v", err),
		}
	}

	if len(utxos) == 0 {
		return SendBitcoinResponse{
			Success: false,
			Message: "사용 가능한 UTXO가 없습니다",
		}
	}

	// 2. 금액 계산 (BTC to satoshi)
	amountSatoshi := int64(request.Amount * 100000000)

	// 3. 입력 선택 및 총 입력 금액 계산
	var totalInput int64
	var selectedUTXOs []UTXO

	// UTXO를 값 순으로 정렬 (큰 것부터)
	sort.Slice(utxos, func(i, j int) bool {
		return utxos[i].Value > utxos[j].Value
	})

	// 필요한 만큼 UTXO 선택
	totalFee := int64(request.FeeSatoshi) // 전체 수수료 (프론트엔드에서 설정)

	// 수수료 분할이 활성화된 경우, 전체 수수료에서 개발자 수수료를 차감한 나머지가 채굴자 수수료
	minerFee := totalFee
	developerFeeNeeded := int64(0)
	if request.EnableFeeSplit {
		developerFeeNeeded = int64(request.DeveloperFeeSatoshi)
		minerFee = totalFee - developerFeeNeeded // 채굴자 수수료 = 전체 수수료 - 개발자 수수료
	}

	totalNeeded := amountSatoshi + totalFee // 전송 금액 + 전체 수수료 (개발자 수수료 포함)
	fmt.Printf("총 필요 금액: %d satoshi (전송: %d + 전체수수료: %d)\n", totalNeeded, amountSatoshi, totalFee)
	if request.EnableFeeSplit {
		fmt.Printf("수수료 분할: 전체수수료 %d = 채굴자수수료 %d + 개발자수수료 %d\n", totalFee, minerFee, developerFeeNeeded)
	}

	for _, utxo := range utxos {
		selectedUTXOs = append(selectedUTXOs, utxo)
		totalInput += utxo.Value

		if totalInput >= totalNeeded {
			break
		}
	}

	// 잔액 확인
	if totalInput < totalNeeded {
		return SendBitcoinResponse{
			Success: false,
			Message: fmt.Sprintf("잔액이 부족합니다. 필요: %d satoshi, 보유: %d satoshi", totalNeeded, totalInput),
		}
	}

	// 4. 거래 생성
	tx := wire.NewMsgTx(wire.TxVersion)

	// 입력 추가
	for _, utxo := range selectedUTXOs {
		txHash, err := hex.DecodeString(utxo.TxID)
		if err != nil {
			return SendBitcoinResponse{
				Success: false,
				Message: fmt.Sprintf("거래 해시 디코딩 실패: %v", err),
			}
		}

		// 바이트 순서 뒤집기 (little-endian)
		for i := 0; i < len(txHash)/2; i++ {
			txHash[i], txHash[len(txHash)-1-i] = txHash[len(txHash)-1-i], txHash[i]
		}

		hash, err := chainhash.NewHash(txHash)
		if err != nil {
			return SendBitcoinResponse{
				Success: false,
				Message: fmt.Sprintf("거래 해시 생성 실패: %v", err),
			}
		}
		outPoint := wire.NewOutPoint(hash, uint32(utxo.Vout))
		txIn := wire.NewTxIn(outPoint, nil, nil)
		tx.AddTxIn(txIn)
	}

	// 받는 주소 파싱
	recipientAddr, err := btcutil.DecodeAddress(request.RecipientAddress, &chaincfg.MainNetParams)
	if err != nil {
		return SendBitcoinResponse{
			Success: false,
			Message: fmt.Sprintf("받는 주소 형식 오류: %v", err),
		}
	}

	// 받는 주소 출력 스크립트 생성
	recipientScript, err := txscript.PayToAddrScript(recipientAddr)
	if err != nil {
		return SendBitcoinResponse{
			Success: false,
			Message: fmt.Sprintf("받는 주소 스크립트 생성 실패: %v", err),
		}
	}

	// 받는 주소 출력 추가
	txOut := wire.NewTxOut(amountSatoshi, recipientScript)
	tx.AddTxOut(txOut)

	// 개발자 수수료 출력 추가 (수수료 분할이 활성화된 경우)
	var developerFeeSatoshi int64 = 0
	if request.EnableFeeSplit && request.DeveloperAddress != "" {
		developerAddr, err := btcutil.DecodeAddress(request.DeveloperAddress, &chaincfg.MainNetParams)
		if err != nil {
			return SendBitcoinResponse{
				Success: false,
				Message: fmt.Sprintf("개발자 주소 형식 오류: %v", err),
			}
		}

		developerScript, err := txscript.PayToAddrScript(developerAddr)
		if err != nil {
			return SendBitcoinResponse{
				Success: false,
				Message: fmt.Sprintf("개발자 주소 스크립트 생성 실패: %v", err),
			}
		}

		developerFeeSatoshi = int64(request.DeveloperFeeSatoshi)
		developerTxOut := wire.NewTxOut(developerFeeSatoshi, developerScript)
		tx.AddTxOut(developerTxOut)

		fmt.Printf("개발자 수수료 출력 추가: %d satoshi → %s\n", developerFeeSatoshi, request.DeveloperAddress)
	}

	// 실제 채굴자 수수료 사용 (전체 수수료에서 개발자 수수료 차감한 값)
	actualMinerFee := minerFee

	// 거스름돈 계산 (개발자 수수료와 채굴자 수수료 모두 차감)
	change := totalInput - amountSatoshi - developerFeeSatoshi - actualMinerFee

	// 거스름돈이 더스트 임계값(546 satoshi)보다 크면 거스름돈 출력 추가
	if change >= 546 {
		changeAddr, err := btcutil.DecodeAddress(request.WalletData.Address, &chaincfg.MainNetParams)
		if err != nil {
			return SendBitcoinResponse{
				Success: false,
				Message: fmt.Sprintf("거스름돈 주소 파싱 실패: %v", err),
			}
		}

		changeScript, err := txscript.PayToAddrScript(changeAddr)
		if err != nil {
			return SendBitcoinResponse{
				Success: false,
				Message: fmt.Sprintf("거스름돈 스크립트 생성 실패: %v", err),
			}
		}

		changeTxOut := wire.NewTxOut(change, changeScript)
		tx.AddTxOut(changeTxOut)
	}

	// 5. 거래 서명
	privKeyWIF, err := btcutil.DecodeWIF(request.WalletData.PrivateKeyWIF)
	if err != nil {
		return SendBitcoinResponse{
			Success: false,
			Message: fmt.Sprintf("개인키 디코딩 실패: %v", err),
		}
	}

	// 각 입력에 대해 서명
	for i, utxo := range selectedUTXOs {
		// 거래 세부정보 조회하여 스크립트 가져오기
		txDetails, err := a.fetchTxDetails(utxo.TxID)
		if err != nil {
			return SendBitcoinResponse{
				Success: false,
				Message: fmt.Sprintf("거래 세부정보 조회 실패: %v", err),
			}
		}

		if utxo.Vout >= len(txDetails.Vout) {
			return SendBitcoinResponse{
				Success: false,
				Message: "잘못된 UTXO 인덱스",
			}
		}

		prevOutScript, err := hex.DecodeString(txDetails.Vout[utxo.Vout].ScriptPubKey)
		if err != nil {
			return SendBitcoinResponse{
				Success: false,
				Message: fmt.Sprintf("이전 출력 스크립트 디코딩 실패: %v", err),
			}
		}

		// P2WPKH 서명 해시 계산 (간단한 fetcher 생성)
		prevOutputFetcher := txscript.NewCannedPrevOutputFetcher(prevOutScript, utxo.Value)
		sigHashes := txscript.NewTxSigHashes(tx, prevOutputFetcher)
		sigHash, err := txscript.CalcWitnessSigHash(prevOutScript, sigHashes, txscript.SigHashAll, tx, i, utxo.Value)
		if err != nil {
			return SendBitcoinResponse{
				Success: false,
				Message: fmt.Sprintf("서명 해시 계산 실패: %v", err),
			}
		}

		// 서명 생성
		signature := ecdsa.Sign(privKeyWIF.PrivKey, sigHash)

		// 서명에 SigHashAll 플래그 추가
		sigWithFlag := append(signature.Serialize(), byte(txscript.SigHashAll))

		// 공개키
		pubKey := privKeyWIF.PrivKey.PubKey().SerializeCompressed()

		// Witness 데이터 설정
		tx.TxIn[i].Witness = wire.TxWitness{sigWithFlag, pubKey}
	}

	// 6. 거래 직렬화
	var buf bytes.Buffer
	if err := tx.Serialize(&buf); err != nil {
		return SendBitcoinResponse{
			Success: false,
			Message: fmt.Sprintf("거래 직렬화 실패: %v", err),
		}
	}

	txHex := hex.EncodeToString(buf.Bytes())

	// 7. 거래 브로드캐스트
	txHash, err := a.broadcastTransaction(txHex)
	if err != nil {
		return SendBitcoinResponse{
			Success: false,
			Message: fmt.Sprintf("거래 브로드캐스트 실패: %v", err),
		}
	}

	return SendBitcoinResponse{
		Success: true,
		Message: "거래가 성공적으로 전송되었습니다",
		TxHash:  txHash,
	}
}
