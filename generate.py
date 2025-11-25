import os
import logging
from typing import Optional

try:
    from llama_cpp import Llama
    LLAMA_CPP_AVAILABLE = True
except ImportError:
    LLAMA_CPP_AVAILABLE = False
    logging.warning("llama-cpp-python not available. Install it for local LLM support.")

class ContentGenerator:
    def __init__(self, model_path: Optional[str] = None):
        self.llm = None
        self.model_path = model_path or self._find_model()
        
        if self.model_path and LLAMA_CPP_AVAILABLE:
            try:
                self.llm = Llama(
                    model_path=self.model_path,
                    n_ctx=2048,
                    n_threads=4,
                    verbose=False
                )
                logging.info(f"Loaded model: {self.model_path}")
            except Exception as e:
                logging.error(f"Failed to load model: {e}")
                self.llm = None
        else:
            logging.warning("No model available. Using fallback generation.")

    def _find_model(self) -> Optional[str]:
        """Find GGUF model in models directory"""
        models_dir = os.path.join(os.path.dirname(__file__), "models")
        if not os.path.exists(models_dir):
            return None
            
        for file in os.listdir(models_dir):
            if file.endswith(('.gguf', '.ggml')):
                return os.path.join(models_dir, file)
        return None

    def load_prompt_template(self, template_name: str = "blog.txt") -> str:
        """Load prompt template from file"""
        template_path = os.path.join(
            os.path.dirname(__file__), 
            "prompt_templates", 
            template_name
        )
        
        try:
            with open(template_path, 'r', encoding='utf-8') as f:
                return f.read()
        except FileNotFoundError:
            logging.error(f"Template not found: {template_path}")
            return "Write a blog article about: {{topic}}"

    def generate_content(self, topic: str) -> str:
        """Generate blog content for given topic"""
        template = self.load_prompt_template()
        prompt = template.replace("{{topic}}", topic)
        
        if self.llm:
            try:
                response = self.llm(
                    prompt,
                    max_tokens=800,
                    temperature=0.7,
                    top_p=0.9,
                    stop=["</s>", "[INST]", "[/INST]"]
                )
                return response['choices'][0]['text'].strip()
            except Exception as e:
                logging.error(f"LLM generation failed: {e}")
                return self._fallback_generation(topic)
        else:
            return self._fallback_generation(topic)

    def _fallback_generation(self, topic: str) -> str:
        """Fallback content generation when LLM is not available"""
        return f"""## OUTLINE
- Introduction to {topic}
- Key aspects and importance
- Practical applications
- Benefits and considerations
- Conclusion and next steps

## ARTICLE

# Understanding {topic}

{topic} is an important subject that deserves our attention and understanding. In today's rapidly evolving world, having knowledge about {topic} can provide significant advantages and insights.

## Key Aspects

When exploring {topic}, several key aspects emerge that are worth considering. These elements form the foundation of our understanding and help us appreciate the complexity and nuance involved.

## Practical Applications

The practical applications of {topic} are numerous and varied. From everyday situations to professional environments, the principles and concepts related to {topic} can be applied in meaningful ways.

## Benefits and Considerations

Understanding {topic} brings several benefits, including improved decision-making, better problem-solving capabilities, and enhanced perspective on related matters. However, it's also important to consider potential challenges and limitations.

## Conclusion

In conclusion, {topic} represents a valuable area of knowledge that can enrich our understanding and provide practical benefits. By continuing to explore and learn about {topic}, we can develop a more comprehensive and nuanced perspective.

*Note: This content was generated using fallback mode. For better quality content, please install a local LLM model.*"""