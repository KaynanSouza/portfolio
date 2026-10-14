package br.com.rocket.projeto_integrador.modules.aiAssistant.factory;

import dev.langchain4j.data.document.Document;
import dev.langchain4j.data.document.DocumentParser;
import dev.langchain4j.data.document.DocumentSplitter;
import dev.langchain4j.data.document.loader.FileSystemDocumentLoader;
import dev.langchain4j.data.document.parser.TextDocumentParser;
import dev.langchain4j.data.document.splitter.DocumentSplitters;
import dev.langchain4j.data.segment.TextSegment;
import dev.langchain4j.model.embedding.EmbeddingModel;
import dev.langchain4j.rag.content.retriever.ContentRetriever;
import dev.langchain4j.rag.content.retriever.EmbeddingStoreContentRetriever;
import dev.langchain4j.store.embedding.EmbeddingStore;

import java.net.URISyntaxException;
import java.net.URL;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;

/**
 * Fábrica responsável por transformar arquivos (um ou vários)
 * em embeddings e criar o ContentRetriever.
 */
public class ContentRetrieverFactory {

    /* ---------- NOVO MÉTODO: suporta N arquivos ---------- */

    /**
     * Constrói um ContentRetriever a partir de uma lista de nomes de arquivos.
     * Todos os segmentos de todos os arquivos são colocados no MESMO EmbeddingStore,
     * permitindo respostas que cruzam informações de fontes distintas.
     */
    public static ContentRetriever createFilesContentRetriever(
            EmbeddingModel embeddingModel,
            EmbeddingStore<TextSegment> embeddingStore,
            List<String> filenames) {

        // 1. Converte cada arquivo em segmentos de texto
        List<TextSegment> allSegments = new ArrayList<>();
        for (String filename : filenames) {
            allSegments.addAll(createTextSegments(filename));
        }

        // 2. Gera embeddings para todos os segmentos
        var embeddings = embeddingModel.embedAll(allSegments).content();

        // 3. Armazena (embedding, segmento) no store
        embeddingStore.addAll(embeddings, allSegments);

        // 4. Cria o retriever único
        return EmbeddingStoreContentRetriever.builder()
                .embeddingStore(embeddingStore)
                .embeddingModel(embeddingModel)
                .maxResults(5)  // pode ajustar
                .minScore(0.5)  // idem
                .build();
    }

    /* ---------- MÉTODO ORIGINAL (single file) — mantido para compatibilidade ---------- */

    public static ContentRetriever createFileContentRetriever(
            EmbeddingModel embeddingModel,
            EmbeddingStore<TextSegment> embeddingStore,
            String filename) {

        return createFilesContentRetriever(embeddingModel, embeddingStore, List.of(filename));
    }

    /* ---------- Utilidades privadas ---------- */

    private static List<TextSegment> createTextSegments(String filename) {
        Path documentPath = toPath(filename);
        DocumentParser documentParser = new TextDocumentParser();
        Document document = FileSystemDocumentLoader.loadDocument(documentPath, documentParser);
        DocumentSplitter splitter = DocumentSplitters.recursive(300, 0);
        return splitter.split(document);
    }

    private static Path toPath(String fileName) {
        try {
            URL url = ContentRetrieverFactory.class.getClassLoader().getResource(fileName);
            assert url != null;
            return Paths.get(url.toURI());
        } catch (URISyntaxException e) {
            throw new RuntimeException(e);
        }
    }
}
