package services

const (
	authorizationHeader = "authorization"
	authorizationBearer = "bearer"
)

// func (u *UserService) authorizeUser(ctx context.Context) (*token.Payload, error) {
// 	md, ok := metadata.FromIncomingContext(ctx)
// 	if !ok {
// 		return nil, fmt.Errorf("missing metadata")
// 	}

// 	values := md.Get(authorizationHeader)
// 	if len(values) == 0 {
// 		return nil, fmt.Errorf("missing authorization header")
// 	}

// 	authHeader := values[0]
// 	fields := strings.Fields(authHeader)
// 	if len(fields) < 2 {
// 		return nil, fmt.Errorf("invalid authorization header format")
// 	}

// 	authType := strings.ToLower(fields[0])
// 	if authType != authorizationBearer {
// 		return nil, fmt.Errorf("unsupported authorization type: %s", authType)
// 	}
// 	tokenMaker, err := token.NewPasetoMaker(u.config.TokenSymmetricKey)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot create token maker: %w", err)
// 	}

// 	accessToken := fields[1]
// 	payload, err := tokenMaker.VerifyToken(accessToken)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid access token: %s", err)
// 	}

// 	return payload, nil
// }
