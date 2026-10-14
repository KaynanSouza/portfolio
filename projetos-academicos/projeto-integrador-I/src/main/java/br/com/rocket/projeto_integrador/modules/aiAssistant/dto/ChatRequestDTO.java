package br.com.rocket.projeto_integrador.modules.aiAssistant.dto;

import java.util.List;

/**
 * NOVO DTO que inclui a lista de arquivos
 * a serem buscados na requisição RAG.
 */
public record ChatRequestDTO(String message, List<String> filenames) { }
