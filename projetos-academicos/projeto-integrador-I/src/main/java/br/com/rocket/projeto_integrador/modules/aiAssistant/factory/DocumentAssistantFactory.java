package br.com.rocket.projeto_integrador.modules.aiAssistant.factory;

import br.com.rocket.projeto_integrador.modules.aiAssistant.utils.DocumentService;
import dev.langchain4j.memory.chat.MessageWindowChatMemory;
import dev.langchain4j.model.chat.ChatLanguageModel;
import dev.langchain4j.rag.content.retriever.ContentRetriever;
import dev.langchain4j.service.AiServices;

/**
 * Encapsula a criação de um service RAG especializado em documentos.
 * Não precisou de mudanças para múltiplos arquivos, pois ele recebe
 * um ContentRetriever já configurado.
 */
public class DocumentAssistantFactory {

    private final DocumentService documentAiService;

    public DocumentAssistantFactory(ChatLanguageModel chatModel, ContentRetriever contentRetriever) {
        this.documentAiService = buildDocumentAiService(chatModel, contentRetriever);
    }

    public String chat(String message) {
        return documentAiService.chat(message);
    }

    /* instancia o proxy gerado por LangChain4j */
    private DocumentService buildDocumentAiService(ChatLanguageModel chatModel, ContentRetriever contentRetriever) {
        return AiServices.builder(DocumentService.class)
                .chatLanguageModel(chatModel)
                .contentRetriever(contentRetriever)
                .chatMemory(MessageWindowChatMemory.withMaxMessages(10))
                .build();
    }
}
