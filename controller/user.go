package controller

// func (c Controller) GetUsers(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	txn := newrelic.FromContext(ctx)
// 	defer txn.StartSegment("GetUsers").End()
// 	startTime := time.Now()

// 	request := new(types.UserBranchRequest)

// 	cfg, ok := r.Context().Value("cfg").(*config.Config)

// 	if !ok {
// 		log.Info().Msg("user is invalid")
// 		log.Error().Interface("req", request).Msg("user is invalid")
// 		response := types.ApiResponse{
// 			Status:  http.StatusBadRequest,
// 			Message: "Data masukan permintaan tidak sesuai",
// 		}
// 		w.WriteHeader(response.Status)
// 		utils.SendHTTPResponse(w, response)
// 		return
// 	}

// 	request.BranchID = cfg.User.AccessID
// 	request.RoleID = cfg.User.RoleID

// 	data, err := c.userSvc.GetUserBranch(r.Context(), *request)
// 	message := "Sukses"
// 	if err != nil || len(data) == 0 {
// 		message = "Pengguna sekolah atau institusi tidak ditemukan"
// 	}
// 	response := types.ApiResponseV2{
// 		Header: types.HeaderResponse{
// 			ProcessTime: time.Since(startTime).Milliseconds(),
// 			Message:     message,
// 		},
// 		Data: data,
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	utils.SendHTTPResponse(w, response)
// }
