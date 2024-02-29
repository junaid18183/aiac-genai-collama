# GenAI
-----
## Run locally 

Build the image
```
nerdctl compose build genai
```

Run it

```
nerdctl run -e OPENAI_API_KEY=$OPENAI_API_KEY --rm  --name genai -p 8086:8086  vivsoft-platform-ui_genai
```

#----
