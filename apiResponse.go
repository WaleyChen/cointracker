package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type APIResponse struct {
	ApiVersion string `json:"apiVersion"`
	RequestId  string `json:"requestId"`
	Context    string `json:"context"`
	Data       Data   `json:"data"`
}

type Data struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Total  int    `json:"total"`
	Items  []Item `json:"items"`
}

type Item struct {
	TransactionId      string             `json:"transactionId"`
	Index              int                `json:"index"`
	MinedInBlockHash   string             `json:"minedInBlockHash"`
	MinedInBlockHeight int                `json:"minedInBlockHeight"`
	Recipients         []Participant      `json:"recipients"`
	Senders            []Participant      `json:"senders"`
	Timestamp          int64              `json:"timestamp"`
	TransactionHash    string             `json:"transactionHash"`
	BlockchainSpecific BlockchainSpecific `json:"blockchainSpecific"`
	Fee                Fee                `json:"fee"`
}

type Participant struct {
	Address string `json:"address"`
	Amount  string `json:"amount"`
}

type BlockchainSpecific struct {
	Locktime int    `json:"locktime"`
	Size     int    `json:"size"`
	VSize    int    `json:"vSize"`
	Version  int    `json:"version"`
	Vin      []Vin  `json:"vin"`
	Vout     []Vout `json:"vout"`
}

type Vin struct {
	Addresses   []string  `json:"addresses"`
	ScriptSig   ScriptSig `json:"scriptSig"`
	Sequence    string    `json:"sequence"`
	Txid        string    `json:"txid"`
	Txinwitness []string  `json:"txinwitness"`
	Value       string    `json:"value"`
	Vout        int       `json:"vout"`
}

type ScriptSig struct {
	Asm  string `json:"asm"`
	Hex  string `json:"hex"`
	Type string `json:"type"`
}

type Vout struct {
	IsSpent      bool         `json:"isSpent"`
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
	Value        string       `json:"value"`
}

type ScriptPubKey struct {
	Addresses []string `json:"addresses"`
	Asm       string   `json:"asm"`
	Hex       string   `json:"hex"`
	Type      string   `json:"type"`
}

type Fee struct {
	Amount string `json:"amount"`
	Unit   string `json:"unit"`
}

func MinInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func calculate_items_on_page(page int, pageLimit int, totalTxs int) int {
	startIndex := page * pageLimit
	maxEndIndex := startIndex + pageLimit
	actualEndIndex := MinInt(maxEndIndex, totalTxs)
	numItemsOnPage := actualEndIndex - startIndex
	return numItemsOnPage
}

func generateFakeItems(page int, pageLimit int, totalTxs int) ([]Item, error) {
	var item Item

	numItems := calculate_items_on_page(page, pageLimit, totalTxs)
	fmt.Println("Num Items: ", numItems)
	items := make([]Item, numItems)

	// Unmarshal the JSON data into the struct
	err := json.Unmarshal([]byte(fakeItem), &item)
	if err != nil {
		return items, errors.New("failed to unmarshal JSON")
	}

	for idx, _ := range items {
		item.TransactionId = uuid.New().String()
		items[idx] = item
	}

	return items, nil
}

const fakeResponse = `
{
    "apiVersion": "2023-04-25",
    "requestId": "65d2fb451c6fef68c2c4081b",
    "context": "yourExampleString",
    "data": {
        "limit": 1,
        "offset": 2,
        "total": 1305156,
        "items": [
            {
                "transactionId": "793b769c67e022961dfc72acff0aeb275d9ea54f6f71450d283e852a79f296a6",
                "index": 30,
                "minedInBlockHash": "0000000000000000000199dcfdf7852c9425a671fffa01e379dc133f5edb6fa5",
                "minedInBlockHeight": 831087,
                "recipients": [
                    {
                        "address": "bc1q33rx2yy9hamzur4ach8ecl7c7kq75lms75arq2",
                        "amount": "1.99984000"
                    },
                    {
                        "address": "bc1qm34lsc65zpw79lxes69zkqmk6ee3ewf0j77s3h",
                        "amount": "3.32049671"
                    },
                    {
                        "address": "bc1qnsdgztu7d4hqrsul9r6nxkasluey9q5yqj56p0",
                        "amount": "0.19812152"
                    },
                    {
                        "address": "18PXSHjJ67K8nvUtHxBiZkTQfa5eQTahZN",
                        "amount": "0.00211462"
                    },
                    {
                        "address": "bc1qxljqqt674qhhz5psgy3g9tvs0wkxgtxq2wzjzp",
                        "amount": "0.00715504"
                    },
                    {
                        "address": "3EDDExUgdCziF1PqG3ie6LtzMb5nZAvPoo",
                        "amount": "0.00180898"
                    },
                    {
                        "address": "bc1q7qmct3nvsjdvde90m3hesc52yp9umgv5k6ze9n",
                        "amount": "0.00112000"
                    },
                    {
                        "address": "3BZa3cUkk9SHXYSeqJhFiNAwhdkbtYp9Et",
                        "amount": "0.06225269"
                    },
                    {
                        "address": "bc1qlqg5seuv23a0w9nkarm24nxrj24z00pj48v3re",
                        "amount": "0.00094060"
                    },
                    {
                        "address": "bc1q9ruv48ye3n2rjlajat9esll5lgge6edal544td",
                        "amount": "32.96699561"
                    },
                    {
                        "address": "bc1qyqcdxtp5xcpg803muuq4xh69mtktww60736xxh",
                        "amount": "0.37649731"
                    }
                ],
                "senders": [
                    {
                        "address": "bc1qm34lsc65zpw79lxes69zkqmk6ee3ewf0j77s3h",
                        "amount": "38.93755558"
                    }
                ],
                "timestamp": 1708325163,
                "transactionHash": "e448efaffbe679e2ef434077de244bf04aab3c833a4680b528a7786315d94a4e",
                "blockchainSpecific": {
                    "locktime": 0,
                    "size": 506,
                    "vSize": 425,
                    "version": 1,
                    "vin": [
                        {
                            "addresses": [
                                "bc1qm34lsc65zpw79lxes69zkqmk6ee3ewf0j77s3h"
                            ],
                            "scriptSig": {
                                "asm": "",
                                "hex": "",
                                "type": "witness_v0_keyhash"
                            },
                            "sequence": "4294967295",
                            "txid": "4820c1c61da2d846dfbd700ad90059aaeebc9492b1479788de97a3b3261f143e",
                            "txinwitness": [
                                "30440220606370bb65bdbdf07fcc7bd85beff2771edd02f56b8263d27ca2f6ffb1cc8f16022061ef2ee12e0fb8e87df0270029299f2a054e8680243c4c418861e4e19d39bb5f01",
                                "02174ee672429ff94304321cdae1fc1e487edf658b34bd1d36da03761658a2bb09"
                            ],
                            "value": "38.93755558",
                            "vout": 1
                        }
                    ],
                    "vout": [
                        {
                            "isSpent": false,
                            "scriptPubKey": {
                                "addresses": [
                                    "bc1q33rx2yy9hamzur4ach8ecl7c7kq75lms75arq2"
                                ],
                                "asm": "0 8c46651085bf762e0ebdc5cf9c7fd8f581ea7f70",
                                "hex": "00148c46651085bf762e0ebdc5cf9c7fd8f581ea7f70",
                                "type": "witness_v0_keyhash"
                            },
                            "value": "1.99984000"
                        },
                        {
                            "isSpent": false,
                            "scriptPubKey": {
                                "addresses": [
                                    "bc1qm34lsc65zpw79lxes69zkqmk6ee3ewf0j77s3h"
                                ],
                                "asm": "0 dc6bf86354105de2fcd9868a2b0376d6731cb92f",
                                "hex": "0014dc6bf86354105de2fcd9868a2b0376d6731cb92f",
                                "type": "witness_v0_keyhash"
                            },
                            "value": "3.32049671"
                        },
                        {
                            "isSpent": false,
                            "scriptPubKey": {
                                "addresses": [
                                    "bc1qnsdgztu7d4hqrsul9r6nxkasluey9q5yqj56p0"
                                ],
                                "asm": "0 9c1a812f9e6d6e01c39f28f5335bb0ff32428284",
                                "hex": "00149c1a812f9e6d6e01c39f28f5335bb0ff32428284",
                                "type": "witness_v0_keyhash"
                            },
                            "value": "0.19812152"
                        },
                        {
                            "isSpent": false,
                            "scriptPubKey": {
                                "addresses": [
                                    "18PXSHjJ67K8nvUtHxBiZkTQfa5eQTahZN"
                                ],
                                "asm": "OP_DUP OP_HASH160 510b703767564d17b85a34a7ebb4cce22aec09b8 OP_EQUALVERIFY OP_CHECKSIG",
                                "hex": "76a914510b703767564d17b85a34a7ebb4cce22aec09b888ac",
                                "type": "pubkeyhash"
                            },
                            "value": "0.00211462"
                        },
                        {
                            "isSpent": false,
                            "scriptPubKey": {
                                "addresses": [
                                    "bc1qxljqqt674qhhz5psgy3g9tvs0wkxgtxq2wzjzp"
                                ],
                                "asm": "0 37e4002f5ea82f715030412282ad907bac642cc0",
                                "hex": "001437e4002f5ea82f715030412282ad907bac642cc0",
                                "type": "witness_v0_keyhash"
                            },
                            "value": "0.00715504"
                        },
                        {
                            "isSpent": false,
                            "scriptPubKey": {
                                "addresses": [
                                    "3EDDExUgdCziF1PqG3ie6LtzMb5nZAvPoo"
                                ],
                                "asm": "OP_HASH160 89576f843e74d1da077ed4e68aa4c2e7cb36accc OP_EQUAL",
                                "hex": "a91489576f843e74d1da077ed4e68aa4c2e7cb36accc87",
                                "type": "scripthash"
                            },
                            "value": "0.00180898"
                        },
                        {
                            "isSpent": false,
                            "scriptPubKey": {
                                "addresses": [
                                    "bc1q7qmct3nvsjdvde90m3hesc52yp9umgv5k6ze9n"
                                ],
                                "asm": "0 f03785c66c849ac6e4afdc6f98628a204bcda194",
                                "hex": "0014f03785c66c849ac6e4afdc6f98628a204bcda194",
                                "type": "witness_v0_keyhash"
                            },
                            "value": "0.00112000"
                        },
                        {
                            "isSpent": false,
                            "scriptPubKey": {
                                "addresses": [
                                    "3BZa3cUkk9SHXYSeqJhFiNAwhdkbtYp9Et"
                                ],
                                "asm": "OP_HASH160 6c48be7f99d14209acb4b91f06ce5819cedc2760 OP_EQUAL",
                                "hex": "a9146c48be7f99d14209acb4b91f06ce5819cedc276087",
                                "type": "scripthash"
                            },
                            "value": "0.06225269"
                        },
                        {
                            "isSpent": false,
                            "scriptPubKey": {
                                "addresses": [
                                    "bc1qlqg5seuv23a0w9nkarm24nxrj24z00pj48v3re"
                                ],
                                "asm": "0 f81148678c547af71676e8f6aaccc392aa27bc32",
                                "hex": "0014f81148678c547af71676e8f6aaccc392aa27bc32",
                                "type": "witness_v0_keyhash"
                            },
                            "value": "0.00094060"
                        },
                        {
                            "isSpent": false,
                            "scriptPubKey": {
                                "addresses": [
                                    "bc1q9ruv48ye3n2rjlajat9esll5lgge6edal544td"
                                ],
                                "asm": "0 28f8ca9c998cd4397fb2eacb987ff4fa119d65bd",
                                "hex": "001428f8ca9c998cd4397fb2eacb987ff4fa119d65bd",
                                "type": "witness_v0_keyhash"
                            },
                            "value": "32.96699561"
                        },
                        {
                            "isSpent": false,
                            "scriptPubKey": {
                                "addresses": [
                                    "bc1qyqcdxtp5xcpg803muuq4xh69mtktww60736xxh"
                                ],
                                "asm": "0 2030d32c34360283be3be701535f45daecb73b4f",
                                "hex": "00142030d32c34360283be3be701535f45daecb73b4f",
                                "type": "witness_v0_keyhash"
                            },
                            "value": "0.37649731"
                        }
                    ]
                },
                "fee": {
                    "amount": "0.00021250",
                    "unit": "BTC"
                }
            }
        ]
    }
}
`

const fakeItem = `
{
	"transactionId": "793b769c67e022961dfc72acff0aeb275d9ea54f6f71450d283e852a79f296a6",
	"index": 30,
	"minedInBlockHash": "0000000000000000000199dcfdf7852c9425a671fffa01e379dc133f5edb6fa5",
	"minedInBlockHeight": 831087,
	"recipients": [
		{
			"address": "bc1q33rx2yy9hamzur4ach8ecl7c7kq75lms75arq2",
			"amount": "1.99984000"
		},
		{
			"address": "bc1qm34lsc65zpw79lxes69zkqmk6ee3ewf0j77s3h",
			"amount": "3.32049671"
		},
		{
			"address": "bc1qnsdgztu7d4hqrsul9r6nxkasluey9q5yqj56p0",
			"amount": "0.19812152"
		},
		{
			"address": "18PXSHjJ67K8nvUtHxBiZkTQfa5eQTahZN",
			"amount": "0.00211462"
		},
		{
			"address": "bc1qxljqqt674qhhz5psgy3g9tvs0wkxgtxq2wzjzp",
			"amount": "0.00715504"
		},
		{
			"address": "3EDDExUgdCziF1PqG3ie6LtzMb5nZAvPoo",
			"amount": "0.00180898"
		},
		{
			"address": "bc1q7qmct3nvsjdvde90m3hesc52yp9umgv5k6ze9n",
			"amount": "0.00112000"
		},
		{
			"address": "3BZa3cUkk9SHXYSeqJhFiNAwhdkbtYp9Et",
			"amount": "0.06225269"
		},
		{
			"address": "bc1qlqg5seuv23a0w9nkarm24nxrj24z00pj48v3re",
			"amount": "0.00094060"
		},
		{
			"address": "bc1q9ruv48ye3n2rjlajat9esll5lgge6edal544td",
			"amount": "32.96699561"
		},
		{
			"address": "bc1qyqcdxtp5xcpg803muuq4xh69mtktww60736xxh",
			"amount": "0.37649731"
		}
	],
	"senders": [
		{
			"address": "bc1qm34lsc65zpw79lxes69zkqmk6ee3ewf0j77s3h",
			"amount": "38.93755558"
		}
	],
	"timestamp": 1708325163,
	"transactionHash": "e448efaffbe679e2ef434077de244bf04aab3c833a4680b528a7786315d94a4e",
	"blockchainSpecific": {
		"locktime": 0,
		"size": 506,
		"vSize": 425,
		"version": 1,
		"vin": [
			{
				"addresses": [
					"bc1qm34lsc65zpw79lxes69zkqmk6ee3ewf0j77s3h"
				],
				"scriptSig": {
					"asm": "",
					"hex": "",
					"type": "witness_v0_keyhash"
				},
				"sequence": "4294967295",
				"txid": "4820c1c61da2d846dfbd700ad90059aaeebc9492b1479788de97a3b3261f143e",
				"txinwitness": [
					"30440220606370bb65bdbdf07fcc7bd85beff2771edd02f56b8263d27ca2f6ffb1cc8f16022061ef2ee12e0fb8e87df0270029299f2a054e8680243c4c418861e4e19d39bb5f01",
					"02174ee672429ff94304321cdae1fc1e487edf658b34bd1d36da03761658a2bb09"
				],
				"value": "38.93755558",
				"vout": 1
			}
		],
		"vout": [
			{
				"isSpent": false,
				"scriptPubKey": {
					"addresses": [
						"bc1q33rx2yy9hamzur4ach8ecl7c7kq75lms75arq2"
					],
					"asm": "0 8c46651085bf762e0ebdc5cf9c7fd8f581ea7f70",
					"hex": "00148c46651085bf762e0ebdc5cf9c7fd8f581ea7f70",
					"type": "witness_v0_keyhash"
				},
				"value": "1.99984000"
			},
			{
				"isSpent": false,
				"scriptPubKey": {
					"addresses": [
						"bc1qm34lsc65zpw79lxes69zkqmk6ee3ewf0j77s3h"
					],
					"asm": "0 dc6bf86354105de2fcd9868a2b0376d6731cb92f",
					"hex": "0014dc6bf86354105de2fcd9868a2b0376d6731cb92f",
					"type": "witness_v0_keyhash"
				},
				"value": "3.32049671"
			},
			{
				"isSpent": false,
				"scriptPubKey": {
					"addresses": [
						"bc1qnsdgztu7d4hqrsul9r6nxkasluey9q5yqj56p0"
					],
					"asm": "0 9c1a812f9e6d6e01c39f28f5335bb0ff32428284",
					"hex": "00149c1a812f9e6d6e01c39f28f5335bb0ff32428284",
					"type": "witness_v0_keyhash"
				},
				"value": "0.19812152"
			},
			{
				"isSpent": false,
				"scriptPubKey": {
					"addresses": [
						"18PXSHjJ67K8nvUtHxBiZkTQfa5eQTahZN"
					],
					"asm": "OP_DUP OP_HASH160 510b703767564d17b85a34a7ebb4cce22aec09b8 OP_EQUALVERIFY OP_CHECKSIG",
					"hex": "76a914510b703767564d17b85a34a7ebb4cce22aec09b888ac",
					"type": "pubkeyhash"
				},
				"value": "0.00211462"
			},
			{
				"isSpent": false,
				"scriptPubKey": {
					"addresses": [
						"bc1qxljqqt674qhhz5psgy3g9tvs0wkxgtxq2wzjzp"
					],
					"asm": "0 37e4002f5ea82f715030412282ad907bac642cc0",
					"hex": "001437e4002f5ea82f715030412282ad907bac642cc0",
					"type": "witness_v0_keyhash"
				},
				"value": "0.00715504"
			},
			{
				"isSpent": false,
				"scriptPubKey": {
					"addresses": [
						"3EDDExUgdCziF1PqG3ie6LtzMb5nZAvPoo"
					],
					"asm": "OP_HASH160 89576f843e74d1da077ed4e68aa4c2e7cb36accc OP_EQUAL",
					"hex": "a91489576f843e74d1da077ed4e68aa4c2e7cb36accc87",
					"type": "scripthash"
				},
				"value": "0.00180898"
			},
			{
				"isSpent": false,
				"scriptPubKey": {
					"addresses": [
						"bc1q7qmct3nvsjdvde90m3hesc52yp9umgv5k6ze9n"
					],
					"asm": "0 f03785c66c849ac6e4afdc6f98628a204bcda194",
					"hex": "0014f03785c66c849ac6e4afdc6f98628a204bcda194",
					"type": "witness_v0_keyhash"
				},
				"value": "0.00112000"
			},
			{
				"isSpent": false,
				"scriptPubKey": {
					"addresses": [
						"3BZa3cUkk9SHXYSeqJhFiNAwhdkbtYp9Et"
					],
					"asm": "OP_HASH160 6c48be7f99d14209acb4b91f06ce5819cedc2760 OP_EQUAL",
					"hex": "a9146c48be7f99d14209acb4b91f06ce5819cedc276087",
					"type": "scripthash"
				},
				"value": "0.06225269"
			},
			{
				"isSpent": false,
				"scriptPubKey": {
					"addresses": [
						"bc1qlqg5seuv23a0w9nkarm24nxrj24z00pj48v3re"
					],
					"asm": "0 f81148678c547af71676e8f6aaccc392aa27bc32",
					"hex": "0014f81148678c547af71676e8f6aaccc392aa27bc32",
					"type": "witness_v0_keyhash"
				},
				"value": "0.00094060"
			},
			{
				"isSpent": false,
				"scriptPubKey": {
					"addresses": [
						"bc1q9ruv48ye3n2rjlajat9esll5lgge6edal544td"
					],
					"asm": "0 28f8ca9c998cd4397fb2eacb987ff4fa119d65bd",
					"hex": "001428f8ca9c998cd4397fb2eacb987ff4fa119d65bd",
					"type": "witness_v0_keyhash"
				},
				"value": "32.96699561"
			},
			{
				"isSpent": false,
				"scriptPubKey": {
					"addresses": [
						"bc1qyqcdxtp5xcpg803muuq4xh69mtktww60736xxh"
					],
					"asm": "0 2030d32c34360283be3be701535f45daecb73b4f",
					"hex": "00142030d32c34360283be3be701535f45daecb73b4f",
					"type": "witness_v0_keyhash"
				},
				"value": "0.37649731"
			}
		]
	},
	"fee": {
		"amount": "0.00021250",
		"unit": "BTC"
	}
}
`
