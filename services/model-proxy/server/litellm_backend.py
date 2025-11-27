


from litellm import completion

async def generate_text(model, prompt, stream=True, **kwargs):
    response = completion(
        model=model,
        prompt=prompt,
        stream=stream,
        **kwargs
    )
    for chunk in response:
        yield chunk["choices"][0]["text"]


