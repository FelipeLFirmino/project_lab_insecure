package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project_lab/internal/middleware"
	"project_lab/internal/models"
	"project_lab/internal/repositories"
	"strconv"
	"strings" // Necessário para o tratamento de erro da versão vulnerável
	"time"
)

type TravelGroupHandler struct {
	repo repositories.TravelGroupRepository
}

func NewTravelGroupHandler(repo repositories.TravelGroupRepository) *TravelGroupHandler {
	return &TravelGroupHandler{repo: repo}
}

// ... CreateGroupHandler e ListGroups (sem alteração de vulnerabilidade)

func (h *TravelGroupHandler) CreateGroupHandler(w http.ResponseWriter, r *http.Request) {

	var req models.TravelGroupCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requisição inválida ou formato JSON incorreto.", http.StatusBadRequest)
		return
	}

	const layout = "2006-01-02"

	startDate, err := time.Parse(layout, req.StartDate)
	if err != nil {
		http.Error(w, "Formato de data de início inválido. Use YYYY-MM-DD.", http.StatusUnprocessableEntity)
		return
	}

	endDate, err := time.Parse(layout, req.EndDate)
	if err != nil {
		http.Error(w, "Formato de data de término inválido. Use YYYY-MM-DD.", http.StatusUnprocessableEntity)
		return
	}

	if req.Name == "" || req.StartDate == "" || req.EndDate == "" {
		http.Error(w, "Nome, data de início e data de término são obrigatórios.", http.StatusUnprocessableEntity)
		return
	}

	if startDate.After(endDate) {
		http.Error(w, "A data de início deve ser anterior ou igual à data de término.", http.StatusUnprocessableEntity)
		return
	}

	userIDValue := r.Context().Value(middleware.UserIDKey)
	creatorID, ok := userIDValue.(int)
	if !ok {
		http.Error(w, "Não autorizado. ID do usuário não encontrado.", http.StatusUnauthorized)
		return
	}

	group := models.TravelGroup{
		Name:        req.Name,
		Description: req.Description,
		StartDate:   startDate,
		EndDate:     endDate,
		CreatorID:   creatorID,
	}

	if err := h.repo.CreateTravelGroup(&group); err != nil {
		fmt.Printf("Erro ao criar grupo no BD: %v\n", err)
		http.Error(w, "Erro interno ao salvar grupo de viagem.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(group); err != nil {
		http.Error(w, "Erro ao serializar resposta.", http.StatusInternalServerError)
		return
	}
}

func (h *TravelGroupHandler) ListGroups(w http.ResponseWriter, r *http.Request) {

	userIdValue := r.Context().Value(middleware.UserIDKey)

	userID, ok := userIdValue.(int)
	if !ok {
		http.Error(w, "Falha na autenticação. ID de usuário não disponível.", http.StatusInternalServerError)
		return
	}

	groups, err := h.repo.ListGroupsByUserId(userID)
	if err != nil {
		fmt.Printf("Erro ao buscar grupos para userID %d: %v\n", userID, err)

		http.Error(w, "Erro interno ao buscar grupos de viagem.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(groups); err != nil {
		http.Error(w, "Erro ao serializar resposta JSON.", http.StatusInternalServerError)
		return
	}
}

// VULNERABILIDADE A01: IDOR INJETADA
// Este handler foi alterado para ignorar o ID do usuário (token JWT) e chamar o repositório
// vulnerável.
func (h *TravelGroupHandler) GetGroupDetailsWithID(w http.ResponseWriter, r *http.Request, groupIDStr string) {

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "ID do grupo inválido. Deve ser um número.", http.StatusBadRequest)
		return
	}

	// CÓDIGO VULNERÁVEL: Removida a lógica de extração e checagem do userID (autorização).

	details, err := h.repo.GetGroupDetails(groupID) // Chamada para o repositório vulnerável (sem userID)
	if err != nil {
		// Ajusta o tratamento de erro para a versão vulnerável
		if strings.Contains(err.Error(), "grupo não encontrado") {
			http.Error(w, "Grupo não encontrado.", http.StatusNotFound)
			return
		}
		fmt.Printf("Erro ao buscar detalhes do grupo %d: %v\n", groupID, err)
		http.Error(w, "Erro interno ao buscar detalhes do grupo.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(details)
}

func (h *TravelGroupHandler) ListGroupMembersHandler(w http.ResponseWriter, r *http.Request, groupIDStr string) {

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "ID do grupo inválido.", http.StatusBadRequest)
		return
	}

	// VULNERABILIDADE A01: Não há checagem de autorização aqui (IDOR).
	// O atacante pode ver os membros de qualquer grupo.
	members, err := h.repo.ListGroupMembers(groupID)
	if err != nil {
		fmt.Printf("Erro ao buscar lista de membros do grupo %d: %v\n", groupID, err)
		http.Error(w, "Erro interno ao buscar membros do grupo.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}

func (h *TravelGroupHandler) ListGroupDestinationsHandler(w http.ResponseWriter, r *http.Request, groupIDStr string) {

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "ID do grupo inválido.", http.StatusBadRequest)
		return
	}

	// VULNERABILIDADE A01: Não há checagem de autorização aqui (IDOR).
	// O atacante pode ver os destinos de qualquer grupo.
	destinations, err := h.repo.ListGroupDestinations(groupID)
	if err != nil {
		fmt.Printf("Erro ao buscar lista de destinos do grupo %d: %v\n", groupID, err)
		http.Error(w, "Erro interno ao buscar destinos do grupo.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(destinations)
}

// VULNERABILIDADE A01: IDOR INJETADA
// Removemos a lógica de extração do userID e a passagem para o repositório.
func (h *TravelGroupHandler) ListGroupVotingsHandler(w http.ResponseWriter, r *http.Request, groupIDStr string) {

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "ID do grupo inválido.", http.StatusBadRequest)
		return
	}

	// CÓDIGO VULNERÁVEL: A extração do userID e a checagem foram removidas.

	votings, err := h.repo.ListGroupVotings(groupID) // Chamada para o repositório vulnerável (sem userID)
	if err != nil {
		fmt.Printf("Erro ao buscar votações do grupo %d: %v\n", groupID, err)
		http.Error(w, "Erro interno ao buscar votações do grupo.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(votings)
}

func (h *TravelGroupHandler) ListGroupExpensesHandler(w http.ResponseWriter, r *http.Request, groupIDStr string) {

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "ID do grupo inválido.", http.StatusBadRequest)
		return
	}

	// VULNERABILIDADE A01: Não há checagem de autorização aqui (IDOR).
	// O atacante pode ver as despesas de qualquer grupo.
	expenses, err := h.repo.ListGroupExpenses(groupID)
	if err != nil {
		fmt.Printf("Erro ao buscar despesas do grupo %d: %v\n", groupID, err)
		http.Error(w, "Erro interno ao buscar despesas do grupo.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

func (h *TravelGroupHandler) CreateDestinationHandler(w http.ResponseWriter, r *http.Request, groupIDStr string) {

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "ID do grupo inválido.", http.StatusBadRequest)
		return
	}

	// VULNERABILIDADE A01: IDOR em Escrita.
	// Nenhuma checagem de autorização de escrita aqui.
	// Qualquer usuário autenticado pode criar um destino em qualquer grupo.

	var req models.DestinationCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requisição inválida (JSON).", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "O nome do destino é obrigatório.", http.StatusUnprocessableEntity)
		return
	}

	destination := models.Destination{
		TravelGroupID: groupID,
		Name:          req.Name,
		Location:      req.Location,
		Description:   req.Description,
	}

	if err := h.repo.CreateDestination(&destination); err != nil {
		fmt.Printf("Erro ao criar destino no BD: %v\n", err)
		http.Error(w, "Erro interno ao salvar destino.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(destination)
}

func (h *TravelGroupHandler) CreateVotingHandler(w http.ResponseWriter, r *http.Request, groupIDStr string) {

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "ID do grupo inválido.", http.StatusBadRequest)
		return
	}

	// VULNERABILIDADE A01: IDOR em Escrita.
	// Nenhuma checagem de autorização de escrita aqui.

	var req models.VotingCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requisição inválida (JSON).", http.StatusBadRequest)
		return
	}

	if req.Question == "" || len(req.Options) < 2 {
		http.Error(w, "A pergunta e pelo menos 2 opções são obrigatórias.", http.StatusUnprocessableEntity)
		return
	}

	optionsJSON, err := json.Marshal(req.Options)
	if err != nil {
		http.Error(w, "Erro ao processar opções da votação.", http.StatusInternalServerError)
		return
	}

	newVotingID, err := h.repo.CreateVoting(groupID, req.Question, string(optionsJSON))
	if err != nil {
		fmt.Printf("Erro ao criar votação no BD: %v\n", err)
		http.Error(w, "Erro interno ao salvar votação.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": newVotingID})
}

func (h *TravelGroupHandler) CreateExpenseHandler(w http.ResponseWriter, r *http.Request, groupIDStr string) {

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "ID do grupo inválido.", http.StatusBadRequest)
		return
	}

	// VULNERABILIDADE A01: IDOR em Escrita.
	// Nenhuma checagem de autorização de escrita aqui.

	var req models.ExpenseCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requisição inválida (JSON).", http.StatusBadRequest)
		return
	}

	if req.Description == "" || req.Amount <= 0 || req.PayerID <= 0 || len(req.ParticipantIDs) == 0 {
		http.Error(w, "Descrição, valor (positivo), pagador e lista de participantes são obrigatórios.", http.StatusUnprocessableEntity)
		return
	}

	expense := models.Expense{
		TravelGroupID:  groupID,
		Description:    req.Description,
		Amount:         req.Amount,
		PayerID:        req.PayerID,
		ParticipantIDs: req.ParticipantIDs,
	}

	if err := h.repo.CreateExpense(&expense); err != nil {
		fmt.Printf("Erro ao criar despesa no BD: %v\n", err)
		http.Error(w, "Erro interno ao salvar despesa e participantes.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(expense)
}
