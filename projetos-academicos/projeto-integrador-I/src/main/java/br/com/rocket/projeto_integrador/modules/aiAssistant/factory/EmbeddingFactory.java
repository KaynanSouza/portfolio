package br.com.rocket.projeto_integrador.modules.aiAssistant.factory;

import dev.langchain4j.model.embedding.AllMiniLmL6V2EmbeddingModel;
import dev.langchain4j.model.embedding.EmbeddingModel;
import dev.langchain4j.store.embedding.EmbeddingStore;
import dev.langchain4j.store.embedding.inmemory.InMemoryEmbeddingStore;

/**
 * Cria o modelo de embeddings e o store correspondente.
 * Para poucas centenas de documentos o store em memória basta.
 */
public class EmbeddingFactory {

    public static EmbeddingModel createEmbeddingModel() {
        return new AllMiniLmL6V2EmbeddingModel();
    }

    public static EmbeddingStore createEmbeddingStore() {
        return new InMemoryEmbeddingStore();
    }
}
