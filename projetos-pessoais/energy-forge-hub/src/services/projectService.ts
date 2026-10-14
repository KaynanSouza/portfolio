import api from "./api";

export const getProjects = async () => {
  try {
    const response = await api.get('/projects');
    return response.data;
  } catch (error) {
    console.error("Erro ao buscar projetos do backend Go:", error);
    throw error;
  }
};

export const getKanbanTasks = async () => {
  try {
    const response = await api.get('/kanban');
    return response.data;
  } catch (error) {
    console.error("Erro ao buscar tarefas do Kanban:", error);
    throw error;
  }
};

export const getDocs = async () => {
  try {
    const response = await api.get('/docs');
    return response.data;
  } catch (error) {
    console.error("Erro ao buscar documentos:", error);
    throw error;
  }
};

export const getTrainings = async () => {
  try {
    const response = await api.get('/trainings');
    return response.data;
  } catch (error) {
    console.error("Erro ao buscar treinamentos:", error);
    throw error;
  }
};

export const createTraining = async (trainingData: any) => {
  try {
    const response = await api.post('/trainings', trainingData);
    return response.data;
  } catch (error) {
    console.error("Erro ao criar treinamento:", error);
    throw error;
  }
};

export const updateTraining = async (id: string, trainingData: any) => {
  try {
    const response = await api.put(`/trainings/${id}`, trainingData);
    return response.data;
  } catch (error) {
    console.error("Erro ao atualizar treinamento:", error);
    throw error;
  }
};

export const deleteTraining = async (id: string) => {
  try {
    const response = await api.delete(`/trainings/${id}`);
    return response.data;
  } catch (error) {
    console.error("Erro ao deletar treinamento:", error);
    throw error;
  }
};