package mailchimp

import (
	"bytes"
	"errors"
	"html/template"
)

var (
	ErrUnknownResource   = errors.New("Unknown resource")
	ErrBadlyFormattedURL = errors.New("Badly formatted URL")
)

type AuthorizedAppParams struct {
	AppID string
}

type AutomationParams struct {
	WorkflowID      string
	WorkflowEmailID string
	SubscriberHash  string
}

type BatchOperationParams struct {
	BatchID string
}

type BatchWebhookParams struct {
	BatchWebhookID string
}

type CampaignFolderParams struct {
	FolderID string
}

type CampaignParams struct {
	CampaignID string
	FeedbackID string
}

type ConversationParams struct {
	ConversationID string
	MessageID      string
}

type EcommerceStoreParams struct {
	StoreID    string
	CartID     string
	LineID     string
	CustomerID string
	OrderID    string
	ProductID  string
	ImageID    string
	VariantID  string
}

type FileManagerFileParams struct {
	FileID string
}

type FileManagerFolderParams struct {
	FolderID string
}

type ListParams struct {
	ListID             string
	ReportID           string
	Month              string
	InterestCategoryID string
	InterestID         string
	SubscriberHash     string
	NoteID             string
	MergeID            string
	SegmentID          string
	WebhookID          string
}

type ReportParams struct {
	CampaignID     string
	ReportID       string
	LinkID         string
	SubscriberHash string
}

type SearchCampaignParams struct{}

type SearchMemberParams struct{}

type TemplateFolderParams struct {
	FolderID string
}

type TemplateParams struct {
	TemplateID string
}

const (
	PostAuthorizedApps                                                       = "/authorized-apps"
	GetAuthorizedApps                                                        = "/authorized-apps"
	GetAuthorizedAppsAppID                                                   = "/authorized-apps/{{.AppID}}"
	GetAutomations                                                           = "/automations"
	GetAutomationsWorkflowID                                                 = "/automations/{{.WorkflowID}}"
	PostAutomationsWorkflowIDActionsPauseAllEmails                           = "/automations/{{.WorkflowID}}/actions/pause-all-emails"
	PostAutomationsWorkflowIDActionsStartAllEmails                           = "/automations/{{.WorkflowID}}/actions/start-all-emails"
	GetAutomationsWorkflowIDEmails                                           = "/automations/{{.WorkflowID}}/emails"
	GetAutomationsWorkflowIDEmailsWorkflowEmailID                            = "/automations/{{.WorkflowID}}/emails/{{.WorkflowEmailID}}"
	PostAutomationsWorkflowIDEmailsWorkflowEmailIDActionsPause               = "/automations/{{.WorkflowID}}/emails/{{.WorkflowEmailID}}/actions/pause"
	PostAutomationsWorkflowIDEmailsWorkflowEmailIDActionsStart               = "/automations/{{.WorkflowID}}/emails/{{.WorkflowEmailID}}/actions/start"
	PostAutomationsWorkflowIDEmailsWorkflowEmailIDQueue                      = "/automations/{{.WorkflowID}}/emails/{{.WorkflowEmailID}}/queue"
	GetAutomationsWorkflowIDEmailsWorkflowEmailIDQueue                       = "/automations/{{.WorkflowID}}/emails/{{.WorkflowEmailID}}/queue"
	GetAutomationsWorkflowIDEmailsWorkflowEmailIDQueueSubscriberHash         = "/automations/{{.WorkflowID}}/emails/{{.WorkflowEmailID}}/queue/{{.SubscriberHash}}"
	PostAutomationsWorkflowIDRemovedSubscribers                              = "/automations/{{.WorkflowID}}/removed-subscribers"
	GetAutomationsWorkflowIDRemovedSubscribers                               = "/automations/{{.WorkflowID}}/removed-subscribers"
	PostBatches                                                              = "/batches"
	GetBatches                                                               = "/batches"
	GetBatchesBatchID                                                        = "/batches/{{.BatchID}}"
	DeleteBatchesBatchID                                                     = "/batches/{{.BatchID}}"
	PostBatchWebhooks                                                        = "/batch-webhooks"
	GetBatchWebhooks                                                         = "/batch-webhooks"
	GetBatchWebhooksBatchWebhookID                                           = "/batch-webhooks/{{.BatchWebhookID}}"
	PatchBatchWebhooksBatchWebhookID                                         = "/batch-webhooks/{{.BatchWebhookID}}"
	DeleteBatchWebhooksBatchWebhookID                                        = "/batch-webhooks/{{.BatchWebhookID}}"
	PostCampaignFolders                                                      = "/campaign-folders"
	GetCampaignFolders                                                       = "/campaign-folders"
	GetCampaignFoldersFolderID                                               = "/campaign-folders/{{.FolderID}}"
	PatchCampaignFoldersFolderID                                             = "/campaign-folders/{{.FolderID}}"
	DeleteCampaignFoldersFolderID                                            = "/campaign-folders/{{.FolderID}}"
	PostCampaigns                                                            = "/campaigns"
	GetCampaigns                                                             = "/campaigns"
	GetCampaignsCampaignID                                                   = "/campaigns/{{.CampaignID}}"
	PatchCampaignsCampaignID                                                 = "/campaigns/{{.CampaignID}}"
	DeleteCampaignsCampaignID                                                = "/campaigns/{{.CampaignID}}"
	PostCampaignsCampaignIDActionsCancelSend                                 = "/campaigns/{{.CampaignID}}/actions/cancel-send"
	PostCampaignsCampaignIDActionsPause                                      = "/campaigns/{{.CampaignID}}/actions/pause"
	PostCampaignsCampaignIDActionsReplicate                                  = "/campaigns/{{.CampaignID}}/actions/replicate"
	PostCampaignsCampaignIDActionsResume                                     = "/campaigns/{{.CampaignID}}/actions/resume"
	PostCampaignsCampaignIDActionsSchedule                                   = "/campaigns/{{.CampaignID}}/actions/schedule"
	PostCampaignsCampaignIDActionsSend                                       = "/campaigns/{{.CampaignID}}/actions/send"
	PostCampaignsCampaignIDActionsTest                                       = "/campaigns/{{.CampaignID}}/actions/test"
	PostCampaignsCampaignIDActionsUnschedule                                 = "/campaigns/{{.CampaignID}}/actions/unschedule"
	GetCampaignsCampaignIDContent                                            = "/campaigns/{{.CampaignID}}/content"
	PutCampaignsCampaignIDContent                                            = "/campaigns/{{.CampaignID}}/content"
	PostCampaignsCampaignIDFeedback                                          = "/campaigns/{{.CampaignID}}/feedback"
	GetCampaignsCampaignIDFeedback                                           = "/campaigns/{{.CampaignID}}/feedback"
	GetCampaignsCampaignIDFeedbackFeedbackID                                 = "/campaigns/{{.CampaignID}}/feedback/{{.FeedbackID}}"
	PatchCampaignsCampaignIDFeedbackFeedbackID                               = "/campaigns/{{.CampaignID}}/feedback/{{.FeedbackID}}"
	DeleteCampaignsCampaignIDFeedbackFeedbackID                              = "/campaigns/{{.CampaignID}}/feedback/{{.FeedbackID}}"
	GetCampaignsCampaignIDSendChecklist                                      = "/campaigns/{{.CampaignID}}/send-checklist"
	GetConversations                                                         = "/conversations"
	GetConversationsConversationID                                           = "/conversations/{{.ConversationID}}"
	PostConversationsConversationIDMessages                                  = "/conversations/{{.ConversationID}}/messages"
	GetConversationsConversationIDMessages                                   = "/conversations/{{.ConversationID}}/messages"
	GetConversationsConversationIDMessagesMessageID                          = "/conversations/{{.ConversationID}}/messages/{{.MessageID}}"
	PostEcommerceStores                                                      = "/ecommerce/stores"
	GetEcommerceStores                                                       = "/ecommerce/stores"
	GetEcommerceStoresStoreID                                                = "/ecommerce/stores/{{.StoreID}}"
	PatchEcommerceStoresStoreID                                              = "/ecommerce/stores/{{.StoreID}}"
	DeleteEcommerceStoresStoreID                                             = "/ecommerce/stores/{{.StoreID}}"
	PostEcommerceStoresStoreIDCarts                                          = "/ecommerce/stores/{{.StoreID}}/carts"
	GetEcommerceStoresStoreIDCarts                                           = "/ecommerce/stores/{{.StoreID}}/carts"
	GetEcommerceStoresStoreIDCartsCartID                                     = "/ecommerce/stores/{{.StoreID}}/carts/{{.CartID}}"
	PatchEcommerceStoresStoreIDCartsCartID                                   = "/ecommerce/stores/{{.StoreID}}/carts/{{.CartID}}"
	DeleteEcommerceStoresStoreIDCartsCartID                                  = "/ecommerce/stores/{{.StoreID}}/carts/{{.CartID}}"
	PostEcommerceStoresStoreIDCartsCartIDLines                               = "/ecommerce/stores/{{.StoreID}}/carts/{{.CartID}}/lines"
	GetEcommerceStoresStoreIDCartsCartIDLines                                = "/ecommerce/stores/{{.StoreID}}/carts/{{.CartID}}/lines"
	GetEcommerceStoresStoreIDCartsCartIDLinesLineID                          = "/ecommerce/stores/{{.StoreID}}/carts/{{.CartID}}/lines/{{.LineID}}"
	PatchEcommerceStoresStoreIDCartsCartIDLinesLineID                        = "/ecommerce/stores/{{.StoreID}}/carts/{{.CartID}}/lines/{{.LineID}}"
	DeleteEcommerceStoresStoreIDCartsCartIDLinesLineID                       = "/ecommerce/stores/{{.StoreID}}/carts/{{.CartID}}/lines/{{.LineID}}"
	PostEcommerceStoresStoreIDCustomers                                      = "/ecommerce/stores/{{.StoreID}}/customers"
	GetEcommerceStoresStoreIDCustomers                                       = "/ecommerce/stores/{{.StoreID}}/customers"
	GetEcommerceStoresStoreIDCustomersCustomerID                             = "/ecommerce/stores/{{.StoreID}}/customers/{{.CustomerID}}"
	PatchEcommerceStoresStoreIDCustomersCustomerID                           = "/ecommerce/stores/{{.StoreID}}/customers/{{.CustomerID}}"
	PutEcommerceStoresStoreIDCustomersCustomerID                             = "/ecommerce/stores/{{.StoreID}}/customers/{{.CustomerID}}"
	DeleteEcommerceStoresStoreIDCustomersCustomerID                          = "/ecommerce/stores/{{.StoreID}}/customers/{{.CustomerID}}"
	PostEcommerceStoresStoreIDOrders                                         = "/ecommerce/stores/{{.StoreID}}/orders"
	GetEcommerceStoresStoreIDOrders                                          = "/ecommerce/stores/{{.StoreID}}/orders"
	GetEcommerceStoresStoreIDOrdersOrderID                                   = "/ecommerce/stores/{{.StoreID}}/orders/{{.OrderID}}"
	PatchEcommerceStoresStoreIDOrdersOrderID                                 = "/ecommerce/stores/{{.StoreID}}/orders/{{.OrderID}}"
	DeleteEcommerceStoresStoreIDOrdersOrderID                                = "/ecommerce/stores/{{.StoreID}}/orders/{{.OrderID}}"
	PostEcommerceStoresStoreIDOrdersOrderIDLines                             = "/ecommerce/stores/{{.StoreID}}/orders/{{.OrderID}}/lines"
	GetEcommerceStoresStoreIDOrdersOrderIDLines                              = "/ecommerce/stores/{{.StoreID}}/orders/{{.OrderID}}/lines"
	GetEcommerceStoresStoreIDOrdersOrderIDLinesLineID                        = "/ecommerce/stores/{{.StoreID}}/orders/{{.OrderID}}/lines/{{.LineID}}"
	PatchEcommerceStoresStoreIDOrdersOrderIDLinesLineID                      = "/ecommerce/stores/{{.StoreID}}/orders/{{.OrderID}}/lines/{{.LineID}}"
	DeleteEcommerceStoresStoreIDOrdersOrderIDLinesLineID                     = "/ecommerce/stores/{{.StoreID}}/orders/{{.OrderID}}/lines/{{.LineID}}"
	PostEcommerceStoresStoreIDProducts                                       = "/ecommerce/stores/{{.StoreID}}/products"
	GetEcommerceStoresStoreIDProducts                                        = "/ecommerce/stores/{{.StoreID}}/products"
	GetEcommerceStoresStoreIDProductsProductID                               = "/ecommerce/stores/{{.StoreID}}/products/{{.ProductID}}"
	PatchEcommerceStoresStoreIDProductsProductID                             = "/ecommerce/stores/{{.StoreID}}/products/{{.ProductID}}"
	DeleteEcommerceStoresStoreIDProductsProductID                            = "/ecommerce/stores/{{.StoreID}}/products/{{.ProductID}}"
	PostEcommerceStoresStoreIDProductsProductIDImages                        = "/ecommerce/stores/{{.StoreID}}/products/{{.ProductID}}/images"
	GetEcommerceStoresStoreIDProductsProductIDImages                         = "/ecommerce/stores/{{.StoreID}}/products/{{.ProductID}}/images"
	GetEcommerceStoresStoreIDProductsProductIDImagesImageID                  = "/ecommerce/stores/{{.StoreID}}/products/{{.ProductID}}/images/{{.ImageID}}"
	PatchEcommerceStoresStoreIDProductsProductIDImagesImageID                = "/ecommerce/stores/{{.StoreID}}/products/{{.ProductID}}/images/{{.ImageID}}"
	DeleteEcommerceStoresStoreIDProductsProductIDImagesImageID               = "/ecommerce/stores/{{.StoreID}}/products/{{.ProductID}}/images/{{.ImageID}}"
	PostEcommerceStoresStoreIDProductsProductIDVariants                      = "/ecommerce/stores/{{.StoreID}}/products/{{.ProductID}}/variants"
	GetEcommerceStoresStoreIDProductsProductIDVariants                       = "/ecommerce/stores/{{.StoreID}}/products/{{.ProductID}}/variants"
	GetEcommerceStoresStoreIDProductsProductIDVariantsVariantID              = "/ecommerce/stores/{{.StoreID}}/products/{{.ProductID}}/variants/{{.VariantID}}"
	PatchEcommerceStoresStoreIDProductsProductIDVariantsVariantID            = "/ecommerce/stores/{{.StoreID}}/products/{{.ProductID}}/variants/{{.VariantID}}"
	PutEcommerceStoresStoreIDProductsProductIDVariantsVariantID              = "/ecommerce/stores/{{.StoreID}}/products/{{.ProductID}}/variants/{{.VariantID}}"
	DeleteEcommerceStoresStoreIDProductsProductIDVariantsVariantID           = "/ecommerce/stores/{{.StoreID}}/products/{{.ProductID}}/variants/{{.VariantID}}"
	PostFileManagerFiles                                                     = "/file-manager/files"
	GetFileManagerFiles                                                      = "/file-manager/files"
	GetFileManagerFilesFileID                                                = "/file-manager/files/{{.FileID}}"
	PatchFileManagerFilesFileID                                              = "/file-manager/files/{{.FileID}}"
	DeleteFileManagerFilesFileID                                             = "/file-manager/files/{{.FileID}}"
	PostFileManagerFolders                                                   = "/file-manager/folders"
	GetFileManagerFolders                                                    = "/file-manager/folders"
	GetFileManagerFoldersFolderID                                            = "/file-manager/folders/{{.FolderID}}"
	PatchFileManagerFoldersFolderID                                          = "/file-manager/folders/{{.FolderID}}"
	DeleteFileManagerFoldersFolderID                                         = "/file-manager/folders/{{.FolderID}}"
	PostLists                                                                = "/lists"
	GetLists                                                                 = "/lists"
	PostListsListID                                                          = "/lists/{{.ListID}}"
	GetListsListID                                                           = "/lists/{{.ListID}}"
	PatchListsListID                                                         = "/lists/{{.ListID}}"
	DeleteListsListID                                                        = "/lists/{{.ListID}}"
	GetListsListIDAbuseReports                                               = "/lists/{{.ListID}}/abuse-reports"
	GetListsListIDAbuseReportsReportID                                       = "/lists/{{.ListID}}/abuse-reports/{{.ReportID}}"
	GetListsListIDActivity                                                   = "/lists/{{.ListID}}/activity"
	GetListsListIDClients                                                    = "/lists/{{.ListID}}/clients"
	GetListsListIDGrowthHistory                                              = "/lists/{{.ListID}}/growth-history"
	GetListsListIDGrowthHistoryMonth                                         = "/lists/{{.ListID}}/growth-history/{{.Month}}"
	PostListsListIDInterestCategories                                        = "/lists/{{.ListID}}/interest-categories"
	GetListsListIDInterestCategories                                         = "/lists/{{.ListID}}/interest-categories"
	GetListsListIDInterestCategoriesInterestCategoryID                       = "/lists/{{.ListID}}/interest-categories/{{.InterestCategoryID}}"
	PatchListsListIDInterestCategoriesInterestCategoryID                     = "/lists/{{.ListID}}/interest-categories/{{.InterestCategoryID}}"
	DeleteListsListIDInterestCategoriesInterestCategoryID                    = "/lists/{{.ListID}}/interest-categories/{{.InterestCategoryID}}"
	PostListsListIDInterestCategoriesInterestCategoryIDInterests             = "/lists/{{.ListID}}/interest-categories/{{.InterestCategoryID}}/interests"
	GetListsListIDInterestCategoriesInterestCategoryIDInterests              = "/lists/{{.ListID}}/interest-categories/{{.InterestCategoryID}}/interests"
	GetListsListIDInterestCategoriesInterestCategoryIDInterestsInterestID    = "/lists/{{.ListID}}/interest-categories/{{.InterestCategoryID}}/interests/{{.InterestID}}"
	PatchListsListIDInterestCategoriesInterestCategoryIDInterestsInterestID  = "/lists/{{.ListID}}/interest-categories/{{.InterestCategoryID}}/interests/{{.InterestID}}"
	DeleteListsListIDInterestCategoriesInterestCategoryIDInterestsInterestID = "/lists/{{.ListID}}/interest-categories/{{.InterestCategoryID}}/interests/{{.InterestID}}"
	GetListsListIDLocations                                                  = "/lists/{{.ListID}}/locations"
	PostListsListIDMembers                                                   = "/lists/{{.ListID}}/members"
	GetListsListIDMembers                                                    = "/lists/{{.ListID}}/members"
	GetListsListIDMembersSubscriberHash                                      = "/lists/{{.ListID}}/members/{{.SubscriberHash}}"
	PatchListsListIDMembersSubscriberHash                                    = "/lists/{{.ListID}}/members/{{.SubscriberHash}}"
	PutListsListIDMembersSubscriberHash                                      = "/lists/{{.ListID}}/members/{{.SubscriberHash}}"
	DeleteListsListIDMembersSubscriberHash                                   = "/lists/{{.ListID}}/members/{{.SubscriberHash}}"
	GetListsListIDMembersSubscriberHashActivity                              = "/lists/{{.ListID}}/members/{{.SubscriberHash}}/activity"
	GetListsListIDMembersSubscriberHashGoals                                 = "/lists/{{.ListID}}/members/{{.SubscriberHash}}/goals"
	PostListsListIDMembersSubscriberHashNotes                                = "/lists/{{.ListID}}/members/{{.SubscriberHash}}/notes"
	GetListsListIDMembersSubscriberHashNotes                                 = "/lists/{{.ListID}}/members/{{.SubscriberHash}}/notes"
	GetListsListIDMembersSubscriberHashNotesNoteID                           = "/lists/{{.ListID}}/members/{{.SubscriberHash}}/notes/{{.NoteID}}"
	PatchListsListIDMembersSubscriberHashNotesNoteID                         = "/lists/{{.ListID}}/members/{{.SubscriberHash}}/notes/{{.NoteID}}"
	DeleteListsListIDMembersSubscriberHashNotesNoteID                        = "/lists/{{.ListID}}/members/{{.SubscriberHash}}/notes/{{.NoteID}}"
	PostListsListIDMergeFields                                               = "/lists/{{.ListID}}/merge-fields"
	GetListsListIDMergeFields                                                = "/lists/{{.ListID}}/merge-fields"
	GetListsListIDMergeFieldsMergeID                                         = "/lists/{{.ListID}}/merge-fields/{{.MergeID}}"
	PatchListsListIDMergeFieldsMergeID                                       = "/lists/{{.ListID}}/merge-fields/{{.MergeID}}"
	DeleteListsListIDMergeFieldsMergeID                                      = "/lists/{{.ListID}}/merge-fields/{{.MergeID}}"
	PostListsListIDSegments                                                  = "/lists/{{.ListID}}/segments"
	GetListsListIDSegments                                                   = "/lists/{{.ListID}}/segments"
	PostListsListIDSegmentsSegmentID                                         = "/lists/{{.ListID}}/segments/{{.SegmentID}}"
	GetListsListIDSegmentsSegmentID                                          = "/lists/{{.ListID}}/segments/{{.SegmentID}}"
	PatchListsListIDSegmentsSegmentID                                        = "/lists/{{.ListID}}/segments/{{.SegmentID}}"
	DeleteListsListIDSegmentsSegmentID                                       = "/lists/{{.ListID}}/segments/{{.SegmentID}}"
	PostListsListIDSegmentsSegmentIDMembers                                  = "/lists/{{.ListID}}/segments/{{.SegmentID}}/members"
	GetListsListIDSegmentsSegmentIDMembers                                   = "/lists/{{.ListID}}/segments/{{.SegmentID}}/members"
	DeleteListsListIDSegmentsSegmentIDMembersSubscriberHash                  = "/lists/{{.ListID}}/segments/{{.SegmentID}}/members/{{.SubscriberHash}}"
	PostListsListIDSignupForms                                               = "/lists/{{.ListID}}/signup-forms"
	GetListsListIDSignupForms                                                = "/lists/{{.ListID}}/signup-forms"
	PostListsListIDWebhooks                                                  = "/lists/{{.ListID}}/webhooks"
	GetListsListIDWebhooks                                                   = "/lists/{{.ListID}}/webhooks"
	GetListsListIDWebhooksWebhookID                                          = "/lists/{{.ListID}}/webhooks/{{.WebhookID}}"
	PatchListsListIDWebhooksWebhookID                                        = "/lists/{{.ListID}}/webhooks/{{.WebhookID}}"
	DeleteListsListIDWebhooksWebhookID                                       = "/lists/{{.ListID}}/webhooks/{{.WebhookID}}"
	GetReports                                                               = "/reports"
	GetReportsCampaignID                                                     = "/reports/{{.CampaignID}}"
	GetReportsCampaignIDAbuseReports                                         = "/reports/{{.CampaignID}}/abuse-reports"
	GetReportsCampaignIDAbuseReportsReportID                                 = "/reports/{{.CampaignID}}/abuse-reports/{{.ReportID}}"
	GetReportsCampaignIDAdvice                                               = "/reports/{{.CampaignID}}/advice"
	GetReportsCampaignIDClickDetails                                         = "/reports/{{.CampaignID}}/click-details"
	GetReportsCampaignIDClickDetailsLinkID                                   = "/reports/{{.CampaignID}}/click-details/{{.LinkID}}"
	GetReportsCampaignIDClickDetailsLinkIDMembers                            = "/reports/{{.CampaignID}}/click-details/{{.LinkID}}/members"
	GetReportsCampaignIDClickDetailsLinkIDMembersSubscriberHash              = "/reports/{{.CampaignID}}/click-details/{{.LinkID}}/members/{{.SubscriberHash}}"
	GetReportsCampaignIDDomainPerformance                                    = "/reports/{{.CampaignID}}/domain-performance"
	GetReportsCampaignIDEepurl                                               = "/reports/{{.CampaignID}}/eepurl"
	GetReportsCampaignIDEmailActivity                                        = "/reports/{{.CampaignID}}/email-activity"
	GetReportsCampaignIDEmailActivitySubscriberHash                          = "/reports/{{.CampaignID}}/email-activity/{{.SubscriberHash}}"
	GetReportsCampaignIDLocations                                            = "/reports/{{.CampaignID}}/locations"
	GetReportsCampaignIDSentTo                                               = "/reports/{{.CampaignID}}/sent-to"
	GetReportsCampaignIDSentToSubscriberHash                                 = "/reports/{{.CampaignID}}/sent-to/{{.SubscriberHash}}"
	GetReportsCampaignIDSubReports                                           = "/reports/{{.CampaignID}}/sub-reports"
	GetReportsCampaignIDUnsubscribed                                         = "/reports/{{.CampaignID}}/unsubscribed"
	GetReportsCampaignIDUnsubscribedSubscriberHash                           = "/reports/{{.CampaignID}}/unsubscribed/{{.SubscriberHash}}"
	GetSearchCampaigns                                                       = "/search-campaigns"
	GetSearchMembers                                                         = "/search-members"
	PostTemplateFolders                                                      = "/template-folders"
	GetTemplateFolders                                                       = "/template-folders"
	GetTemplateFoldersFolderID                                               = "/template-folders/{{.FolderID}}"
	PatchTemplateFoldersFolderID                                             = "/template-folders/{{.FolderID}}"
	DeleteTemplateFoldersFolderID                                            = "/template-folders/{{.FolderID}}"
	PostTemplates                                                            = "/templates"
	GetTemplates                                                             = "/templates"
	GetTemplatesTemplateID                                                   = "/templates/{{.TemplateID}}"
	PatchTemplatesTemplateID                                                 = "/templates/{{.TemplateID}}"
	DeleteTemplatesTemplateID                                                = "/templates/{{.TemplateID}}"
	GetTemplatesTemplateIDDefaultContent                                     = "/templates/{{.TemplateID}}/default-content"
)

var endpointURLs = []string{
	PostAuthorizedApps,
	GetAuthorizedApps,
	GetAuthorizedAppsAppID,
	GetAutomations,
	GetAutomationsWorkflowID,
	PostAutomationsWorkflowIDActionsPauseAllEmails,
	PostAutomationsWorkflowIDActionsStartAllEmails,
	GetAutomationsWorkflowIDEmails,
	GetAutomationsWorkflowIDEmailsWorkflowEmailID,
	PostAutomationsWorkflowIDEmailsWorkflowEmailIDActionsPause,
	PostAutomationsWorkflowIDEmailsWorkflowEmailIDActionsStart,
	PostAutomationsWorkflowIDEmailsWorkflowEmailIDQueue,
	GetAutomationsWorkflowIDEmailsWorkflowEmailIDQueue,
	GetAutomationsWorkflowIDEmailsWorkflowEmailIDQueueSubscriberHash,
	PostAutomationsWorkflowIDRemovedSubscribers,
	GetAutomationsWorkflowIDRemovedSubscribers,
	PostBatches,
	GetBatches,
	GetBatchesBatchID,
	DeleteBatchesBatchID,
	PostBatchWebhooks,
	GetBatchWebhooks,
	GetBatchWebhooksBatchWebhookID,
	PatchBatchWebhooksBatchWebhookID,
	DeleteBatchWebhooksBatchWebhookID,
	PostCampaignFolders,
	GetCampaignFolders,
	GetCampaignFoldersFolderID,
	PatchCampaignFoldersFolderID,
	DeleteCampaignFoldersFolderID,
	PostCampaigns,
	GetCampaigns,
	GetCampaignsCampaignID,
	PatchCampaignsCampaignID,
	DeleteCampaignsCampaignID,
	PostCampaignsCampaignIDActionsCancelSend,
	PostCampaignsCampaignIDActionsPause,
	PostCampaignsCampaignIDActionsReplicate,
	PostCampaignsCampaignIDActionsResume,
	PostCampaignsCampaignIDActionsSchedule,
	PostCampaignsCampaignIDActionsSend,
	PostCampaignsCampaignIDActionsTest,
	PostCampaignsCampaignIDActionsUnschedule,
	GetCampaignsCampaignIDContent,
	PutCampaignsCampaignIDContent,
	PostCampaignsCampaignIDFeedback,
	GetCampaignsCampaignIDFeedback,
	GetCampaignsCampaignIDFeedbackFeedbackID,
	PatchCampaignsCampaignIDFeedbackFeedbackID,
	DeleteCampaignsCampaignIDFeedbackFeedbackID,
	GetCampaignsCampaignIDSendChecklist,
	GetConversations,
	GetConversationsConversationID,
	PostConversationsConversationIDMessages,
	GetConversationsConversationIDMessages,
	GetConversationsConversationIDMessagesMessageID,
	PostEcommerceStores,
	GetEcommerceStores,
	GetEcommerceStoresStoreID,
	PatchEcommerceStoresStoreID,
	DeleteEcommerceStoresStoreID,
	PostEcommerceStoresStoreIDCarts,
	GetEcommerceStoresStoreIDCarts,
	GetEcommerceStoresStoreIDCartsCartID,
	PatchEcommerceStoresStoreIDCartsCartID,
	DeleteEcommerceStoresStoreIDCartsCartID,
	PostEcommerceStoresStoreIDCartsCartIDLines,
	GetEcommerceStoresStoreIDCartsCartIDLines,
	GetEcommerceStoresStoreIDCartsCartIDLinesLineID,
	PatchEcommerceStoresStoreIDCartsCartIDLinesLineID,
	DeleteEcommerceStoresStoreIDCartsCartIDLinesLineID,
	PostEcommerceStoresStoreIDCustomers,
	GetEcommerceStoresStoreIDCustomers,
	GetEcommerceStoresStoreIDCustomersCustomerID,
	PatchEcommerceStoresStoreIDCustomersCustomerID,
	PutEcommerceStoresStoreIDCustomersCustomerID,
	DeleteEcommerceStoresStoreIDCustomersCustomerID,
	PostEcommerceStoresStoreIDOrders,
	GetEcommerceStoresStoreIDOrders,
	GetEcommerceStoresStoreIDOrdersOrderID,
	PatchEcommerceStoresStoreIDOrdersOrderID,
	DeleteEcommerceStoresStoreIDOrdersOrderID,
	PostEcommerceStoresStoreIDOrdersOrderIDLines,
	GetEcommerceStoresStoreIDOrdersOrderIDLines,
	GetEcommerceStoresStoreIDOrdersOrderIDLinesLineID,
	PatchEcommerceStoresStoreIDOrdersOrderIDLinesLineID,
	DeleteEcommerceStoresStoreIDOrdersOrderIDLinesLineID,
	PostEcommerceStoresStoreIDProducts,
	GetEcommerceStoresStoreIDProducts,
	GetEcommerceStoresStoreIDProductsProductID,
	PatchEcommerceStoresStoreIDProductsProductID,
	DeleteEcommerceStoresStoreIDProductsProductID,
	PostEcommerceStoresStoreIDProductsProductIDImages,
	GetEcommerceStoresStoreIDProductsProductIDImages,
	GetEcommerceStoresStoreIDProductsProductIDImagesImageID,
	PatchEcommerceStoresStoreIDProductsProductIDImagesImageID,
	DeleteEcommerceStoresStoreIDProductsProductIDImagesImageID,
	PostEcommerceStoresStoreIDProductsProductIDVariants,
	GetEcommerceStoresStoreIDProductsProductIDVariants,
	GetEcommerceStoresStoreIDProductsProductIDVariantsVariantID,
	PatchEcommerceStoresStoreIDProductsProductIDVariantsVariantID,
	PutEcommerceStoresStoreIDProductsProductIDVariantsVariantID,
	DeleteEcommerceStoresStoreIDProductsProductIDVariantsVariantID,
	PostFileManagerFiles,
	GetFileManagerFiles,
	GetFileManagerFilesFileID,
	PatchFileManagerFilesFileID,
	DeleteFileManagerFilesFileID,
	PostFileManagerFolders,
	GetFileManagerFolders,
	GetFileManagerFoldersFolderID,
	PatchFileManagerFoldersFolderID,
	DeleteFileManagerFoldersFolderID,
	PostLists,
	GetLists,
	PostListsListID,
	GetListsListID,
	PatchListsListID,
	DeleteListsListID,
	GetListsListIDAbuseReports,
	GetListsListIDAbuseReportsReportID,
	GetListsListIDActivity,
	GetListsListIDClients,
	GetListsListIDGrowthHistory,
	GetListsListIDGrowthHistoryMonth,
	PostListsListIDInterestCategories,
	GetListsListIDInterestCategories,
	GetListsListIDInterestCategoriesInterestCategoryID,
	PatchListsListIDInterestCategoriesInterestCategoryID,
	DeleteListsListIDInterestCategoriesInterestCategoryID,
	PostListsListIDInterestCategoriesInterestCategoryIDInterests,
	GetListsListIDInterestCategoriesInterestCategoryIDInterests,
	GetListsListIDInterestCategoriesInterestCategoryIDInterestsInterestID,
	PatchListsListIDInterestCategoriesInterestCategoryIDInterestsInterestID,
	DeleteListsListIDInterestCategoriesInterestCategoryIDInterestsInterestID,
	GetListsListIDLocations,
	PostListsListIDMembers,
	GetListsListIDMembers,
	GetListsListIDMembersSubscriberHash,
	PatchListsListIDMembersSubscriberHash,
	PutListsListIDMembersSubscriberHash,
	DeleteListsListIDMembersSubscriberHash,
	GetListsListIDMembersSubscriberHashActivity,
	GetListsListIDMembersSubscriberHashGoals,
	PostListsListIDMembersSubscriberHashNotes,
	GetListsListIDMembersSubscriberHashNotes,
	GetListsListIDMembersSubscriberHashNotesNoteID,
	PatchListsListIDMembersSubscriberHashNotesNoteID,
	DeleteListsListIDMembersSubscriberHashNotesNoteID,
	PostListsListIDMergeFields,
	GetListsListIDMergeFields,
	GetListsListIDMergeFieldsMergeID,
	PatchListsListIDMergeFieldsMergeID,
	DeleteListsListIDMergeFieldsMergeID,
	PostListsListIDSegments,
	GetListsListIDSegments,
	PostListsListIDSegmentsSegmentID,
	GetListsListIDSegmentsSegmentID,
	PatchListsListIDSegmentsSegmentID,
	DeleteListsListIDSegmentsSegmentID,
	PostListsListIDSegmentsSegmentIDMembers,
	GetListsListIDSegmentsSegmentIDMembers,
	DeleteListsListIDSegmentsSegmentIDMembersSubscriberHash,
	PostListsListIDSignupForms,
	GetListsListIDSignupForms,
	PostListsListIDWebhooks,
	GetListsListIDWebhooks,
	GetListsListIDWebhooksWebhookID,
	PatchListsListIDWebhooksWebhookID,
	DeleteListsListIDWebhooksWebhookID,
	GetReports,
	GetReportsCampaignID,
	GetReportsCampaignIDAbuseReports,
	GetReportsCampaignIDAbuseReportsReportID,
	GetReportsCampaignIDAdvice,
	GetReportsCampaignIDClickDetails,
	GetReportsCampaignIDClickDetailsLinkID,
	GetReportsCampaignIDClickDetailsLinkIDMembers,
	GetReportsCampaignIDClickDetailsLinkIDMembersSubscriberHash,
	GetReportsCampaignIDDomainPerformance,
	GetReportsCampaignIDEepurl,
	GetReportsCampaignIDEmailActivity,
	GetReportsCampaignIDEmailActivitySubscriberHash,
	GetReportsCampaignIDLocations,
	GetReportsCampaignIDSentTo,
	GetReportsCampaignIDSentToSubscriberHash,
	GetReportsCampaignIDSubReports,
	GetReportsCampaignIDUnsubscribed,
	GetReportsCampaignIDUnsubscribedSubscriberHash,
	GetSearchCampaigns,
	GetSearchMembers,
	PostTemplateFolders,
	GetTemplateFolders,
	GetTemplateFoldersFolderID,
	PatchTemplateFoldersFolderID,
	DeleteTemplateFoldersFolderID,
	PostTemplates,
	GetTemplates,
	GetTemplatesTemplateID,
	PatchTemplatesTemplateID,
	DeleteTemplatesTemplateID,
	GetTemplatesTemplateIDDefaultContent,
}
var endpointTemplates = map[string]*template.Template{}

func init() {
	// initialize url templates
	for _, u := range endpointURLs {
		tmpl, err := template.New(u).Parse(u)
		if err != nil {
			panic(err)
		}
		endpointTemplates[u] = tmpl
	}
}

func formatUrl(url string, params interface{}) (string, error) {
	tmpl, ok := endpointTemplates[url]
	if !ok {
		return "", ErrUnknownResource
	}
	if params == nil {
		return url, nil
	}
	b := bytes.Buffer{}
	err := tmpl.Execute(&b, params)
	if err != nil {
		return "", ErrBadlyFormattedURL
	}
	return b.String(), nil
}
