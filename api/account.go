package api

import (
	"database/sql"
	"net/http"
	

	db "github.com/Somraj2929/simple-banking/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// Server struct
type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}


type getAccountListRequest struct {
	PageID int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

type updateAccountURI struct {
    ID int64 `uri:"id" binding:"required,min=1"`
}

type updateAccountBody struct {
    Balance int64 `json:"balance" binding:"required,min=0"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name(){
				case "unique_violation", "foreign_key_violation":
					ctx.JSON(http.StatusForbidden, errorResponse(err))
					return

				// case "foreign_key_violation":
				// 	ctx.JSON(http.StatusNotFound, errorResponse(err))
				// 	return	
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(200, account)
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}



func (server *Server) listAccounts(ctx *gin.Context) {
	var req getAccountListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(200, accounts)
}

func (server *Server) updateAccount(ctx *gin.Context) {
  
    var uriReq updateAccountURI
    if err := ctx.ShouldBindUri(&uriReq); err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }


    var bodyReq updateAccountBody
    if err := ctx.ShouldBindJSON(&bodyReq); err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }


    arg := db.UpdateAccountParams{
        ID:      uriReq.ID,
        Balance: bodyReq.Balance,
    }

    account, err := server.store.UpdateAccount(ctx, arg)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, account)
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(200, gin.H{"status": "account deleted"})
}