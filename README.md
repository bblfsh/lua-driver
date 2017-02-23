# lua-driver

### Example input program

```lua
local function lt(a, b)
        return a < b
end
```

### Example output

```json
{
  "ast": {
    "Stmts": [
      {
        "LocalDecl": false,
        "LocalFunc": true,
        "Targets": [
          {
            "Value": "lt"
          }
        ],
        "Values": [
          {
            "Params": [
              "a",
              "b"
            ],
            "IsVariadic": false,
            "Source": "",
            "Block": [
              {
                "Items": [
                  {
                    "Op": "OpLessThan",
                    "Left": {
                      "Value": "a"
                    },
                    "Right": {
                      "Value": "b"
                    }
                  }
                ]
              }
            ]
          }
        ]
      }
    ]
  },
  "status": "ok",
  "errors": [],
  "language": "lua",
  "language_version": "1.0.0",
  "driver": "lua:1.0.0"
}
```
