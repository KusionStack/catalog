schema Inference:
    """ Inference is a module schema consisting of model, framework and so on

    Attributes
    ----------
    model: str, default is Undefined, required.
        The model name to be used for inference. 
    framework: "Ollama" | "KubeRay", default is Undefined, required.
        The framework or environment in which the model operates.
    system: str, default is "".
        The system message, which will be set in the template.
    template: str, default is "".
        The full prompt template, which will be sent to the model.
    top_k: int, default is 40.
        A higher value (e.g. 100) will give more diverse answers, while a lower value (e.g. 10) will be more conservative.
    top_p: float, default is 0.9.
        A higher value (e.g. 0.9) will give more diverse answers, while a lower value (e.g. 0.5) will be more conservative.
    num_predict: int, default is 128.
        Maximum number of tokens to predict when generating text.
    num_ctx: int, default is 2048.
        The size of the context window used to generate the next token.
    
    Examples
    --------
    import inference.v1.infer

    accessories: {
        "inference@v0.1.0": infer.Inference {
            model: "llama3"
            framework: "Ollama"

            system: "You are Mario from super mario bros, acting as an assistant."
            template: "{{ if .System }}<|im_start|>system {{ .System }}<|im_end|> {{ end }}{{ if .Prompt }}<|im_start|>user {{ .Prompt }}<|im_end|> {{ end }}<|im_start|>assistant"

            top_k: 40
            top_p: 0.9
            temperature: 0.8

            num_predict: 128
            num_ctx: 2048
        }
    }
    """
    model: str
    framework: "Ollama" | "KubeRay"
    system?: str = ""
    template?: str = ""
    top_k?: int = 40
    top_p?: float = 0.9
    temperature?: float = 0.8
    num_predict?: int = 128
    num_ctx?: int = 2048

    check:
        0 < top_k if top_k, "top_k must be more than 0"
        0 < top_p <= 1 if top_p, "top_p must be greater than 0 and less than or equal to 1"
        0 < temperature if temperature, "temperature must be more than 0"
        -2 <= num_predict if num_predict, "num_predict must be greater than or equal to -2"
        0 < num_ctx if num_ctx, "num_ctx must be greater than 0"

