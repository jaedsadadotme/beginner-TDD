package demo

import (
	"fmt"
	"net/http"
)

type Storer interface {
	Save(Demo) error
}

type Handler struct {
	channel string
	store   Storer
}

type Context interface {
	Demo() (Demo, error)
	JSON(int, interface{})
}

func (h *Handler) Demo(c Context) {
	order, err := c.Demo()
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if order.Datas != h.channel {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": fmt.Sprintf("%s is not accepted", order.Datas),
		})
		return
	}

	if err := h.store.Save(order); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"message": "ok",
	})
}
