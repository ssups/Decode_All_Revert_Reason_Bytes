package main

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var qs = []*survey.Question{
	{
		Name:     "revertReason",
		Prompt:   &survey.Input{Message: "wirte revert reason"},
		Validate: survey.Required,
	},
}

func main() {
	answers := struct {
		RevertReason string
	}{}

	for {
		fmt.Println("")
		err := survey.Ask(qs, &answers)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		revertReasonHex, _ := strings.CutPrefix(answers.RevertReason, "0x") // remove prefix 0x
		if len(revertReasonHex) < 8 {
			fmt.Println("wrong input")
			continue
		}

		selector := revertReasonHex[:8]
		switch selector {
		case "08c379a0": // selector for Error(string) -> normal revert
			onlyStringHex := revertReasonHex[8+64+64:] // remove selector, offset, length
			res, err := hexToAscii(onlyStringHex)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(res)

		case "4e487b71": // selector for Panic(uint256) -> panic revert
			onlyUintByte := revertReasonHex[8:] // remove selector
			var startInd int
			for ; startInd < len(onlyUintByte); startInd++ {
				if onlyUintByte[startInd:startInd+1] != "0" {
					break
				}
			}
			panicCode := "0x" + onlyUintByte[startInd:]
			fmt.Println("reverted with panic code: ", panicCode)

		default: // other selector is custom error
			fmt.Println("this is custom error, abi needed to decode")
		}
	}
}

func hexToAscii(hex string) (string, error) {
	hex, _ = strings.CutPrefix(hex, "0x") // remove prefix 0x
	var res string
	for i := 0; i < len(hex)/2; i++ {
		charHex := hex[i*2 : (i+1)*2]
		if charHex == "00" {
			break
		}
		code, err := hexutil.DecodeUint64("0x" + hex[i*2:(i+1)*2])
		if err != nil {
			fmt.Println("here")
			return "", err
		}
		res += string(rune(code))
	}
	return res, nil
}
