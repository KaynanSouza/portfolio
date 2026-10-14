import { useState, useEffect } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Clock, Users, Plus, Loader2, X, Edit, Trash2 } from "lucide-react";
import { getTrainings, createTraining, updateTraining, deleteTraining } from "@/services/projectService";
import { toast } from "sonner";

type TrainingItem = {
    id?: string;
    title: string;
    category: string;
    duration: string;
    enrolled: number;
    status: string;
};

const defaultForm = { title: "", category: "", duration: "", status: "disponível" };

const Training = () => {
    const [trainings, setTrainings] = useState<TrainingItem[]>([]);
    const [loading, setLoading] = useState(true);
    const [isFormOpen, setIsFormOpen] = useState(false);
    const [submitting, setSubmitting] = useState(false);
    const [editingId, setEditingId] = useState<string | null>(null);

    const [formData, setFormData] = useState(defaultForm);

    useEffect(() => {
        fetchTrainings();
    }, []);

    const fetchTrainings = async () => {
        try {
            setLoading(true);
            const data = await getTrainings();
            setTrainings(data || []);
        } catch (error) {
            toast.error("Erro ao carregar os treinamentos.");
        } finally {
            setLoading(false);
        }
    };

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
        const { name, value } = e.target;
        setFormData((prev) => ({ ...prev, [name]: value }));
    };

    const resetForm = () => {
        setFormData(defaultForm);
        setEditingId(null);
        setIsFormOpen(false);
    };

    const handleEdit = (t: TrainingItem) => {
        setFormData({ title: t.title, category: t.category, duration: t.duration, status: t.status });
        setEditingId(t.id || null);
        setIsFormOpen(true);
        window.scrollTo({ top: 0, behavior: "smooth" });
    };

    const handleDelete = async (id: string | undefined) => {
        if (!id) return;
        if (!window.confirm("Tem certeza que deseja apagar este treinamento?")) return;

        try {
            await deleteTraining(id);
            toast.success("Treinamento apagado com sucesso!");
            fetchTrainings();
        } catch (error) {
            toast.error("Erro ao apagar o treinamento.");
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!formData.title || !formData.category || !formData.duration) {
            toast.warning("Preencha todos os campos obrigatórios.");
            return;
        }

        try {
            setSubmitting(true);

            if (editingId) {
                // Atualiza existente
                await updateTraining(editingId, formData);
                toast.success("Treinamento atualizado com sucesso!");
            } else {
                // Cria novo
                await createTraining({ ...formData, enrolled: 0 });
                toast.success("Treinamento cadastrado com sucesso!");
            }

            fetchTrainings();
            resetForm();
        } catch (error) {
            toast.error("Erro ao salvar o treinamento.");
        } finally {
            setSubmitting(false);
        }
    };

    return (
        <div className="space-y-6 max-w-6xl mx-auto">
            <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                <div>
                    <h1 className="text-2xl font-bold">Treinamentos</h1>
                    <p className="text-sm text-muted-foreground">Capacitação técnica do time</p>
                </div>
                <Button
                    onClick={isFormOpen ? resetForm : () => setIsFormOpen(true)}
                    className="flex items-center gap-2"
                    variant={isFormOpen ? "outline" : "default"}
                >
                    {isFormOpen ? <X className="h-4 w-4" /> : <Plus className="h-4 w-4" />}
                    {isFormOpen ? "Cancelar" : "Novo Treinamento"}
                </Button>
            </div>

            {/* Formulário (Usado para Criar e Editar) */}
            {isFormOpen && (
                <Card className="glass-card animate-in fade-in slide-in-from-top-4 border-primary/20 bg-primary/5">
                    <CardHeader>
                        <CardTitle className="text-lg">
                            {editingId ? "Editar Treinamento" : "Cadastrar Novo Treinamento"}
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <form onSubmit={handleSubmit} className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4 items-end">
                            <div className="space-y-1 lg:col-span-2">
                                <label className="text-xs font-medium">Título do Curso</label>
                                <Input name="title" value={formData.title} onChange={handleInputChange} placeholder="Ex: NR-10 Básico" />
                            </div>
                            <div className="space-y-1">
                                <label className="text-xs font-medium">Categoria</label>
                                <Input name="category" value={formData.category} onChange={handleInputChange} placeholder="Ex: Segurança" />
                            </div>
                            <div className="space-y-1">
                                <label className="text-xs font-medium">Duração</label>
                                <Input name="duration" value={formData.duration} onChange={handleInputChange} placeholder="Ex: 8h" />
                            </div>
                            <div className="space-y-1">
                                <label className="text-xs font-medium">Status</label>
                                <select
                                    name="status"
                                    value={formData.status}
                                    onChange={handleInputChange}
                                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                                >
                                    <option value="disponível">Disponível</option>
                                    <option value="em breve">Em Breve</option>
                                </select>
                            </div>
                            <div className="lg:col-span-5 flex justify-end mt-2">
                                <Button type="submit" disabled={submitting} className="w-full md:w-auto">
                                    {submitting ? <Loader2 className="h-4 w-4 animate-spin mr-2" /> : null}
                                    {editingId ? "Salvar Alterações" : "Salvar Treinamento"}
                                </Button>
                            </div>
                        </form>
                    </CardContent>
                </Card>
            )}

            {/* Listagem */}
            {loading ? (
                <div className="flex justify-center py-12">
                    <Loader2 className="h-8 w-8 animate-spin text-primary" />
                </div>
            ) : trainings.length === 0 ? (
                <div className="text-center py-12 text-muted-foreground border-2 border-dashed rounded-lg">
                    Nenhum treinamento cadastrado no banco de dados.
                </div>
            ) : (
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                    {trainings.map((t) => (
                        <Card key={t.id} className="glass-card hover:shadow-md transition-shadow relative group">
                            <CardHeader className="pb-2">
                                <div className="flex items-start justify-between">
                                    <CardTitle className="text-base font-medium pr-2 leading-tight">{t.title}</CardTitle>

                                    {/* Container de Ações visível no hover */}
                                    <div className="flex items-center gap-1 opacity-100 sm:opacity-0 sm:group-hover:opacity-100 transition-opacity bg-background/80 p-1 rounded-md absolute top-3 right-3 shadow-sm border">
                                        <Button variant="ghost" size="icon" className="h-7 w-7 text-muted-foreground hover:text-primary" onClick={() => handleEdit(t)}>
                                            <Edit className="h-3.5 w-3.5" />
                                        </Button>
                                        <Button variant="ghost" size="icon" className="h-7 w-7 text-muted-foreground hover:text-destructive" onClick={() => handleDelete(t.id)}>
                                            <Trash2 className="h-3.5 w-3.5" />
                                        </Button>
                                    </div>
                                </div>

                                <div className="flex items-center gap-2 mt-2">
                                    <Badge
                                        variant="outline"
                                        className={
                                            t.status === "disponível"
                                                ? "bg-success/15 text-success border-success/30"
                                                : "bg-muted text-muted-foreground"
                                        }
                                    >
                                        {t.status}
                                    </Badge>
                                    <Badge variant="secondary" className="text-xs">{t.category}</Badge>
                                </div>
                            </CardHeader>
                            <CardContent>
                                <div className="flex items-center gap-4 text-xs text-muted-foreground mt-2 border-t pt-3">
                  <span className="flex items-center gap-1 font-medium">
                    <Clock className="h-3.5 w-3.5" />
                      {t.duration}
                  </span>
                                    <span className="flex items-center gap-1 font-medium">
                    <Users className="h-3.5 w-3.5" />
                                        {t.enrolled} inscritos
                  </span>
                                </div>
                            </CardContent>
                        </Card>
                    ))}
                </div>
            )}
        </div>
    );
};

export default Training;