package br.com.rocket.projeto_integrador.modules.aiAssistant.factory;

import dev.langchain4j.model.chat.ChatLanguageModel;
import dev.langchain4j.model.huggingface.HuggingFaceChatModel;
import dev.langchain4j.model.openai.OpenAiChatModel;

import java.time.Duration;

/**
 * Fábrica de modelos de linguagem.
 * Mantivemos o código original, adicionando apenas comentários.
 */
public class AiAssistantFactory {

    /**
     * Cria modelo hospedado no Hugging Face (exige token).
     */
    public static ChatLanguageModel createHuggingFace(String accessToken) {
        return HuggingFaceChatModel.builder()
                .accessToken(accessToken)
                .modelId("gemma-2-9b-it")
                .timeout(Duration.ofSeconds(1000))
                .build();
    }

    /**
     * Cria modelo local compatível com a API OpenAI (ex.: LM Studio, Ollama).
     */
    public static ChatLanguageModel createLocalChatModel() {
        return OpenAiChatModel.builder()
                .baseUrl("http://localhost:1234/v1")
                .apiKey("ignore")
                .logRequests(true)
                .timeout(Duration.ofSeconds(1000))
                .build();
    }
}
