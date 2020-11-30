package sqs

type Handler struct {
	Router *Router
}

func (h *Handler) HandleMessage(message SNSMessage) (bool, error) {
	b, err := h.Router.Handle(message)
	if err != nil {
		return false, err
	}
	return b, nil
}
