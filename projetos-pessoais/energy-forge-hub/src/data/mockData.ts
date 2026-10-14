// ============================================================
// Mock Data — replace with API calls (Microsoft Graph, GitLab)
// ============================================================

export interface Project {
  id: string;
  name: string;
  description: string;
  lastCommit: string;
  lastCommitDate: string;
  branch: string;
  ciStatus: "success" | "failed" | "running";
  url: string;
  web_url: string;
}

export interface Task {
  id: string;
  title: string;
  priority: "alta" | "média" | "baixa";
  dueDate: string;
  assignee: string;
  status: "todo" | "in_progress" | "done";
}

export interface DocItem {
  id: string;
  title: string;
  filename: string;
  size: string;
  type: "pdf" | "bin" | "zip" | "txt" | "docx";
  content: string;
}

export interface DocCategory {
  id: string;
  label: string;
  icon: string;
  items: DocItem[];
}

export const summaryCards = {
  activeProjects: 12,
  pendingTasks: 27,
  serversOnline: 8,
  serversTotal: 10,
};


export const tasks: Task[] = [
  { id: "1", title: "Revisar PR do módulo SCADA", priority: "alta", dueDate: "2026-03-10", assignee: "Carlos M.", status: "todo" },
  { id: "2", title: "Atualizar documentação IEC 61850", priority: "média", dueDate: "2026-03-12", assignee: "Ana P.", status: "todo" },
  { id: "3", title: "Corrigir bug no parser de firmware", priority: "alta", dueDate: "2026-03-11", assignee: "João S.", status: "in_progress" },
  { id: "4", title: "Deploy do EnergyDash v3.2", priority: "média", dueDate: "2026-03-13", assignee: "Maria L.", status: "in_progress" },
  { id: "5", title: "Teste de integração OMS", priority: "baixa", dueDate: "2026-03-15", assignee: "Pedro R.", status: "done" },
  { id: "6", title: "Configurar pipeline CI/CD", priority: "alta", dueDate: "2026-03-09", assignee: "Carlos M.", status: "todo" },
  { id: "7", title: "Migrar banco de dados legado", priority: "média", dueDate: "2026-03-14", assignee: "Ana P.", status: "in_progress" },
  { id: "8", title: "Calibrar sensores RTU", priority: "baixa", dueDate: "2026-03-16", assignee: "João S.", status: "done" },
];

export const docCategories: DocCategory[] = [
  {
    id: "equipamentos",
    label: "Equipamentos",
    icon: "Cpu",
    items: [
      { 
        id: "e1", 
        title: "Guia de Instalação RTU-3000", 
        filename: "RTU-3000_Installation_Guide.pdf",
        size: "2.4 MB",
        type: "pdf",
        content: "Conteúdo mockado para o Guia de Instalação RTU-3000..." 
      },
      { 
        id: "e2", 
        title: "Especificação Relé SEL-751", 
        filename: "SEL-751_Specs.pdf",
        size: "1.8 MB",
        type: "pdf",
        content: "Conteúdo mockado para a Especificação do Relé SEL-751..." 
      },
      { 
        id: "e3", 
        title: "Manual de Operação CLP-500", 
        filename: "CLP-500_Operation_Manual.pdf",
        size: "3.1 MB",
        type: "pdf",
        content: "Conteúdo mockado para o Manual de Operação CLP-500..." 
      },
    ],
  },
  {
    id: "manuais",
    label: "Manuais",
    icon: "BookOpen",
    items: [
      { 
        id: "m1", 
        title: "Padrão de Código Python v2.0", 
        filename: "Python_Coding_Standard_v2.pdf",
        size: "850 KB",
        type: "pdf",
        content: "Diretrizes para desenvolvimento em Python nos projetos de engenharia..." 
      },
      { 
        id: "m2", 
        title: "Procedimento de Deploy Global", 
        filename: "Global_Deploy_Procedure.docx",
        size: "1.2 MB",
        type: "docx",
        content: "Checklist e fluxos para deploy em ambientes de dev, staging e production..." 
      },
      { 
        id: "m3", 
        title: "Guia de Onboarding Engenharia", 
        filename: "Engineering_Onboarding_Guide.pdf",
        size: "4.5 MB",
        type: "pdf",
        content: "Bem-vindo à equipe de engenharia. Aqui estão os primeiros passos..." 
      },
    ],
  },
  {
    id: "firmwares",
    label: "Firmwares",
    icon: "Microchip",
    items: [
      { 
        id: "f1", 
        title: "Firmware RTU-3000 v2.4.1", 
        filename: "RTU-3000_v2.4.1.bin",
        size: "12.5 MB",
        type: "bin",
        content: "Arquivo binário do Firmware RTU-3000 v2.4.1" 
      },
      { 
        id: "f2", 
        title: "Hotfix SCADA Gateway v1.1.2", 
        filename: "SCADA_Gateway_v1.1.2_hotfix.zip",
        size: "45.2 MB",
        type: "zip",
        content: "Pacote de hotfix para o SCADA Gateway" 
      },
      { 
        id: "f3", 
        title: "Stack Protocolo IEC 61850 v3.0", 
        filename: "IEC61850_Stack_v3.0.zip",
        size: "8.7 MB",
        type: "zip",
        content: "Implementação da stack do protocolo IEC 61850" 
      },
    ],
  },
];

