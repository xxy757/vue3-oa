package service

import (
	"oa-saas/internal/model"
	"oa-saas/internal/pkg/utils"
	"oa-saas/internal/repository"
)

type NoticeService struct {
	noticeRepo *repository.NoticeRepo
	userRepo   *repository.UserRepo
}

func NewNoticeService(noticeRepo *repository.NoticeRepo, userRepo *repository.UserRepo) *NoticeService {
	return &NoticeService{
		noticeRepo: noticeRepo,
		userRepo:   userRepo,
	}
}

type NoticeItem struct {
	model.Notice
	PublisherName string `json:"publisherName"`
	IsRead        bool   `json:"isRead"`
}

func (s *NoticeService) List(tenantID uint, noticeType int, keyword string, page, pageSize int, userID uint) ([]NoticeItem, int64, error) {
	notices, total, err := s.noticeRepo.ListByTenant(tenantID, noticeType, keyword, page, pageSize)
	if err != nil {
		return nil, 0, utils.ErrInternal("获取公告列表失败")
	}
	var list []NoticeItem
	for _, n := range notices {
		item := NoticeItem{Notice: n}
		if user, err := s.userRepo.GetByID(n.PublisherID, tenantID); err == nil {
			item.PublisherName = user.Nickname
		}
		isRead, _ := s.noticeRepo.HasRead(n.ID, userID, tenantID)
		item.IsRead = isRead
		list = append(list, item)
	}
	return list, total, nil
}

func (s *NoticeService) Detail(id, tenantID, userID uint) (*model.Notice, error) {
	notice, err := s.noticeRepo.GetByID(id, tenantID)
	if err != nil {
		return nil, utils.ErrNotFound("公告不存在")
	}
	read := &model.NoticeRead{}
	s.noticeRepo.GetOrCreateRead(uint(id), userID, tenantID, read)
	return notice, nil
}

func (s *NoticeService) Create(tenantID, publisherID uint, title, content string, noticeType int8, summary, cover string, isTop, status int8) (*model.Notice, error) {
	notice := model.Notice{
		TenantID:    tenantID,
		Title:       title,
		Content:     content,
		Type:        noticeType,
		Summary:     summary,
		Cover:       cover,
		IsTop:       isTop,
		Status:      status,
		PublisherID: publisherID,
	}
	if err := s.noticeRepo.Create(&notice); err != nil {
		return nil, utils.ErrInternal("发布公告失败")
	}
	return &notice, nil
}

func (s *NoticeService) UnreadCount(userID, tenantID uint) (int64, error) {
	count, err := s.noticeRepo.CountUnread(userID, tenantID)
	if err != nil {
		return 0, utils.ErrInternal("查询未读数失败")
	}
	return count, nil
}

func (s *NoticeService) MarkRead(noticeID, userID, tenantID uint) error {
	read := &model.NoticeRead{}
	return s.noticeRepo.GetOrCreateRead(noticeID, userID, tenantID, read)
}
