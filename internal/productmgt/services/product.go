package services

import (
	"context"
	"fmt"

	"github.com/glu/shopvui/golibs/database"
	"github.com/glu/shopvui/idl/pb"
	"github.com/glu/shopvui/internal/productmgt/entities"
	"github.com/glu/shopvui/internal/productmgt/repositories"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProductService struct {
	//config util.Config
	pb.UnimplementedProductServiceServer
	DB database.Ext

	ProductRepo interface {
		InsertProduct(ctx context.Context, db database.Ext, e *entities.Product) (*entities.Product, error)
		GetProductByID(ctx context.Context, db database.Ext, id pgtype.Text) (*entities.Product, error)
		ListAllProducts(ctx context.Context, db database.Ext) ([]*entities.Product, error)
	}
}

func NewProductService(db database.Ext) pb.ProductServiceServer {
	return &ProductService{
		DB:          db,
		ProductRepo: new(repositories.ProductRepo),
	}
}

func toInsertProduct(req *pb.InsertProductRequest) *entities.Product {
	return &entities.Product{
		ID:            database.Text(uuid.NewString()),
		Sku:           database.Text(req.Product.Sku),
		Name:          database.Text(req.Product.Name),
		Description:   database.Text(req.Product.Description),
		RegularPrice:  database.Int4(req.Product.RegularPrice),
		DiscountPrice: database.Int4(req.Product.DiscountPrice),
		Quantity:      database.Int4(req.Product.Quantity),
		Taxable:       database.Bool(req.Product.Taxable),
	}
}

func (s *ProductService) InsertProduct(ctx context.Context, req *pb.InsertProductRequest) (*pb.InsertProductResponse, error) {
	fmt.Println("come InsertProduct ne`")
	product := toInsertProduct(req)
	fmt.Println("check toInsertProduct()")
	p, err := s.ProductRepo.InsertProduct(ctx, s.DB, product)
	fmt.Println("check  s.ProductRepo.InsertProduct(ctx, s.DB, product)")
	if err != nil {
		return nil, err
	}
	return &pb.InsertProductResponse{
		Product: toProductResponse(p),
	}, nil
}

func toProductResponse(p *entities.Product) *pb.Product {
	return &pb.Product{
		Id:            p.ID.String,
		Sku:           p.ID.String,
		Name:          p.Name.String,
		Description:   p.Description.String,
		RegularPrice:  p.RegularPrice.Int,
		DiscountPrice: p.DiscountPrice.Int,
		Quantity:      p.Quantity.Int,
		Taxable:       p.Taxable.Bool,
		CreatedAt:     timestamppb.New(p.CreatedAt.Time),
		UpdatedAt:     timestamppb.New(p.UpdatedAt.Time),
	}
}

func (s *ProductService) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	if req.ProductId == "" {
		return nil, fmt.Errorf("ProductId is required")
	}
	p, err := s.ProductRepo.GetProductByID(ctx, s.DB, database.Text(req.ProductId))
	if err != nil {
		return nil, fmt.Errorf("can't get product: %v", err)
	}
	product := toProductResponse(p)
	return &pb.ListProductsResponse{
		Product: product,
	}, nil

}

func toProductsRsp(p []*entities.Product) []*pb.Product {
	rsp := make([]*pb.Product, len(p))
	for i := 0; i < len(p); i++ {
		rsp = append(rsp, toProductResponse(p[i]))
	}
	return rsp
}

func (s *ProductService) ListAllProducts(ctx context.Context, req *pb.ListAllProductsRequest) (*pb.ListAllProductsResponse, error) {
	//todo: check permissions

	p, err := s.ProductRepo.ListAllProducts(ctx, s.DB)
	if err != nil {
		return nil, fmt.Errorf("cannot list all products: %v", err)
	}

	products := toProductsRsp(p)
	return &pb.ListAllProductsResponse{
		Products: products,
	}, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	return &pb.DeleteProductResponse{
		Message: "success",
	}, nil
}
