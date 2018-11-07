package tmpl

//CollectionTemplate represents a 2.1 schema postman template
const CollectionTemplate = `{
	"info": {
		"_postman_id": "00208b92-7c63-48a2-ae29-4c26bf59c41c",
		"name": "{{.Info.Name}}",
		{{with .Info.Description }} "description": "{{.Info.Description}}", {{end}}
		"schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
	},
	"item": [ {{ $itemCount := len .Items }}	{{ range $iindex, $item :=  .Items }}
		{
			"name": "{{$item.Name}}",
			"event": [{{ $itemEventCount := len $item.Event }}	{{ range $ieindex, $itemevent :=  $item.Event }}
				{
					"listen": "{{$itemevent.Listen}}",
					"script": {
						"id": "{{ uuid }}",
						"exec": [ {{ $itemExecCount := len $itemevent.Script.Exec }} {{ range $iexeindex, $exec :=  $itemevent.Script.Exec }}
							{{ raw $exec}} {{endItem $iexeindex $itemExecCount }} {{end}}
						],
						"type": "{{$itemevent.Script.Type}}"
					}
				} {{endItem $ieindex $itemEventCount }}	 {{end}}
			],
			"request": {
				"method": "{{$item.Request.Method}}",
				"header": [ {{ $headCount := len  $item.Request.Headers  }}	{{ range $heindex, $header :=  $item.Request.Headers }}
					{
						"key": "{{$header.Key}}",
						"value": {{raw $header.Value}}
					}  {{endItem $heindex $headCount }} {{ end}}
				],
				"body": {
					"mode": "{{$item.Request.Body.Mode}}",
					"raw": {{ raw $item.Request.Body.Raw}}
				},
				"url": {
					"raw": {{ raw $item.Request.URL.Raw}},
					"host": [
						{{ $hostCount := len  $item.Request.URL.Host  }}	{{ range $hindex, $host :=  $item.Request.URL.Host }}"{{$host}}"{{endItem $hindex $hostCount }}
						{{ end}}
					],
					"path": [ {{ $pathCount := len  $item.Request.URL.Path  }}	{{ range $pindex, $path :=  $item.Request.URL.Path }}"{{$path}}"{{endItem $pindex $pathCount }}{{ end}}
					]
				}
			},
			"response": []
			} {{endItem $iindex $itemCount }}	
			{{end}}
	],
	"event": [ {{ $eventCount := len .Event }}	{{ range $eindex, $event :=  .Event }}
		{
			"listen": "{{$event.Listen}}",
			"script": {
				"id": "{{ uuid }}",
				"type": "{{$event.Script.Type}}",
				"exec": [ {{ $itemExecCount := len $event.Script.Exec }} {{ range $iexeindex, $exec :=  $event.Script.Exec }}
							{{ raw $exec}} {{endItem $iexeindex $itemExecCount }} {{end}}
						]
			}
		} {{endItem $eindex $eventCount }}	
	{{ end}}
	],
	"variable": [	{{ $varCount := len .Variables }} {{ range $index, $element :=  .Variables }}
			{
				"id": "{{ uuid }}",
				"key": "{{$element.Key}}",
				"value": "{{$element.Value}}",
				"type": "{{$element.Type}}"
			} {{endItem $index $varCount }}	
		{{ end}}
	]
}`
