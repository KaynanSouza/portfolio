package br.com.rocket.projeto_integrador.modules.aiAssistant.controller;

import br.com.rocket.projeto_integrador.modules.aiAssistant.dto.ChatRequestDTO;
import br.com.rocket.projeto_integrador.modules.aiAssistant.factory.AiAssistantFactory;
import br.com.rocket.projeto_integrador.modules.aiAssistant.factory.ContentRetrieverFactory;
import br.com.rocket.projeto_integrador.modules.aiAssistant.factory.DocumentAssistantFactory;
import br.com.rocket.projeto_integrador.modules.aiAssistant.factory.EmbeddingFactory;
import dev.langchain4j.model.chat.ChatLanguageModel;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

/**
 * REST controller que expõe o endpoint /api/chat
 * capaz de receber 1 ou N arquivos para consulta.
 */
@RestController
@RequestMapping("/api/chat")
public class AiAssistantController {

    @Value("${langchain.huggingface.accessToken}")
    private String token;

    /**
     * Endpoint que aceita um JSON no formato:
     * {
     *   "message": "Pergunta do usuário",
     *   "filenames": ["movies.txt", "directors.txt"]
     * }
     *
     * Se 'filenames' vier vazio ou nulo, usa um default.
     */
    @PostMapping
    public ResponseEntity<String> chat(@RequestBody ChatRequestDTO request) {

        /* ----------- 1. Modelos e stores ----------- */
        ChatLanguageModel chatModel = AiAssistantFactory.createLocalChatModel();
        var embeddingModel  = EmbeddingFactory.createEmbeddingModel();
        var embeddingStore  = EmbeddingFactory.createEmbeddingStore();

        /* ----------- 2. Constroi o retriever ----------- */
        List<String> files = request.filenames() == null || request.filenames().isEmpty()
                ? List.of("movies.txt")            // fallback
                : request.filenames();

        var contentRetriever = ContentRetrieverFactory.createFilesContentRetriever(
                embeddingModel,
                embeddingStore,
                files);

        /* ----------- 3. Instancia serviço RAG e responde ----------- */
        var documentAssistant = new DocumentAssistantFactory(chatModel, contentRetriever);
        String answer = documentAssistant.chat(request.message());

        // 4. Resposta HTTP 200 (OK) com o texto gerado
        return ResponseEntity.ok(answer);
    }
}
