/*
 * Copyright (C) distroy
 */

package copybenchstruct1

type ItemCardData struct {
	CardSet     *ItemCardSet `json:"card_set"`
	ItemCards   []*ItemCard  `json:"item_cards"`
	CardVersion *string      `json:"card_version"`
}

type ItemCardSet struct {
	LayoutId      *uint32        `json:"layout_id"`
	ElementToggle *ElementToggle `json:"element_toggle"`
	CardSetName   *string        `json:"card_set_name"`
}

type ElementToggle struct {
	AtcButton      *bool `json:"atc_button"`
	FeedbackButton *bool `json:"feedback_button"`
}

type ItemCard struct {
	ItemCardDisplayedAsset *ItemCardDisplayedAsset `json:"item_card_displayed_asset"`
	ItemData               *ItemData               `json:"item_data"`
}

type ItemCardDisplayedAsset struct {
	Name                  *string                `json:"name"`
	Image                 *string                `json:"image"`
	Images                []string               `json:"images"`
	ShopLocation          *string                `json:"shop_location"`
	HighlightVideos       []*VideoInfo           `json:"highlight_videos"`
	ItemCardMask          *ItemCardMask          `json:"item_card_mask"`
	IconInImage           *IconInImage           `json:"icon_in_image"`
	ImageOverlay          *ImageOverlay          `json:"image_overlay"`
	SellerFlag            *SellerFlag            `json:"seller_flag"`
	PromotionLabelList    []*PromotionLabel      `json:"promotion_label_list"`
	IconInPrice           *IconInPrice           `json:"icon_in_price"`
	DiscountTag           *DiscountTag           `json:"discount_tag"`
	VoucherAppliedIcon    *VoucherAppliedIcon    `json:"voucher_applied_icon"`
	Rating                *Rating                `json:"rating"`
	SoldCount             *SoldCount             `json:"sold_count"`
	EstimatedDeliveryTime *EstimatedDeliveryTime `json:"estimated_delivery_time"`
	FreeBadge             *FreeBadge             `json:"free_badge"`
	FindSimilar           *FindSimilar           `json:"find_similar"`
	LongImage             *string                `json:"long_image"`
	LongImages            []string               `json:"long_images"`
	TopProductBadge       *TopProductBadge       `json:"top_product_badge"`
	ItemCardSpecs         *ItemCardSpecs         `json:"item_card_specs"`
	DisplayPrice          *DisplayPrice          `json:"display_price"`
	CtaButton             *CTAButton             `json:"cta_button"`
}

type ItemData struct {
	Itemid                   *int64                    `json:"itemid"`
	Shopid                   *int64                    `json:"shopid"`
	IsAdult                  *bool                     `json:"is_adult"`
	NeedKyc                  *bool                     `json:"need_kyc"`
	AdultAgeThreshold        *uint32                   `json:"adult_age_threshold"`
	ItemCardDisplayPrice     *ItemCardDisplayPrice     `json:"item_card_display_price"`
	ItemCardDisplaySoldCount *ItemCardDisplaySoldCount `json:"item_card_display_sold_count"`
	IsSoldOut                *bool                     `json:"is_sold_out"`
	IsPreview                *bool                     `json:"is_preview"`
	LabelIds                 []int64                   `json:"label_ids"`
	Catid                    *int32                    `json:"catid"`
	ShopeeVerified           *bool                     `json:"shopee_verified"`
	LikedCount               *int32                    `json:"liked_count"`
	ItemType                 *uint32                   `json:"item_type"`
	ReferenceItemId          *string                   `json:"reference_item_id"`
	ShopData                 *ShopData                 `json:"shop_data"`
	Status                   *int32                    `json:"status"`
	Ctime                    *int32                    `json:"ctime"`
	Flag                     *int32                    `json:"flag"`
	ItemStatus               *string                   `json:"item_status"`
	GlobalCat                *CategoryPath             `json:"global_cat"`
	OverlayImages            []*OverlayImage           `json:"overlay_images"`
	ItemRating               *ItemRating               `json:"item_rating"`
	GlobalBrand              *Brand                    `json:"global_brand"`
	PlatformVoucher          *PlatformVoucher          `json:"platform_voucher"`
	VideoInfoList            []*VideoInfo              `json:"video_info_list"`
	TierVariations           []*TierVariation          `json:"tier_variations"`
	CanUseBundleDeal         *bool                     `json:"can_use_bundle_deal"`
}

type ItemCardDisplaySoldCount struct {
	HistoricalSoldCount     *uint64 `json:"historical_sold_count"`
	MonthlySoldCount        *uint64 `json:"monthly_sold_count"`
	HistoricalSoldCountText *string `json:"historical_sold_count_text"`
	MonthlySoldCountText    *string `json:"monthly_sold_count_text"`
}

type ItemCardDisplayPrice struct {
	ItemId                                *uint64 `json:"item_id"`
	ModelId                               *uint64 `json:"model_id"`
	PromotionType                         *uint32 `json:"promotion_type"`
	PromotionId                           *uint64 `json:"promotion_id"`
	Price                                 *int64  `json:"price"`
	StrikethroughPrice                    *int64  `json:"strikethrough_price"`
	Discount                              *int32  `json:"discount"`
	HiddenPriceDisplayText                *string `json:"hidden_price_display_text"`
	RecommendedShopVoucherPromotionId     *int64  `json:"recommended_shop_voucher_promotion_id"`
	RecommendedPlatformVoucherPromotionId *int64  `json:"recommended_platform_voucher_promotion_id"`
}

type ShopData struct {
	ShopName     *string `json:"shop_name"`
	ShopIcon     *string `json:"shop_icon"`
	ShopLocation *string `json:"shop_location"`
}

type CategoryPath struct {
	Catid []int32 `json:"catid"`
}

type OverlayImage struct {
	OverlayImage *string  `json:"overlay_image"`
	OverlayIds   []string `json:"overlay_ids"`
}

type ItemRating struct {
	RatingStar  *float64 `json:"rating_star"`
	RatingCount []int32  `json:"rating_count"`
}

type Brand struct {
	BrandId     *uint32 `json:"brand_id"`
	DisplayName *string `json:"display_name"`
}

type PlatformVoucher struct {
	IsNuv     *bool    `json:"is_nuv"`
	VoucherId []uint64 `json:"voucher_id"`
}

type TierProperty struct {
	Name  *string `json:"name"`
	Image *string `json:"image"`
	Color *string `json:"color"`
}

type TierVariation struct {
	Name       *string         `json:"name"`
	Options    []string        `json:"options"`
	Images     []string        `json:"images"`
	Properties []*TierProperty `json:"properties"`
	Type       *uint32         `json:"type"`
}

type TopProductBadge struct {
	Rank *uint32 `json:"rank"`
}

type FindSimilar struct {
	ButtonText *string `json:"button_text"`
}

type FreeBadge struct {
	BadgeText *string `json:"badge_text"`
}

type EstimatedDeliveryTime struct {
	EstimatedDeliveryTimeText     *string `json:"estimated_delivery_time_text"`
	MinEstimatedDeliveryTimestamp *int64  `json:"min_estimated_delivery_timestamp"`
	MaxEstimatedDeliveryTimestamp *int64  `json:"max_estimated_delivery_timestamp"`
	EdtType                       *uint32 `json:"edt_type"`
	Distance                      *string `json:"distance"`
	EdtIcon                       *Image  `json:"edt_icon"`
	FailureReason                 *int32  `json:"failure_reason"`
}

type SoldCount struct {
	Text *string `json:"text"`
}

type Rating struct {
	RatingText *string `json:"rating_text"`
	Icon       *Image  `json:"icon"`
	RatingType *uint32 `json:"rating_type"`
}

type VoucherAppliedIcon struct {
	CanDisplay *bool `json:"can_display"`
}

type DiscountTag struct {
	DiscountText *string `json:"discount_text"`
}

type IconInPrice struct {
	Name *string `json:"name"`
	Icon *Image  `json:"icon"`
}

type PromotionLabel struct {
	Type  *string              `json:"type"`
	Style *string              `json:"style"`
	Specs *PromotionLabelSpecs `json:"specs"`
	Data  *PromotionLabelData  `json:"data"`
}

type PromotionLabelData struct {
	Text            *string `json:"text"`
	TextLeft        *string `json:"text_left"`
	TextRight       *string `json:"text_right"`
	ProgressPercent *uint32 `json:"progress_percent"`
}

type PromotionLabelSpecs struct {
	TextColor       *string     `json:"text_color"`
	BorderColor     *string     `json:"border_color"`
	BackgroundColor *string     `json:"background_color"`
	PrimaryColor    *string     `json:"primary_color"`
	SecondaryColor  *string     `json:"secondary_color"`
	Image           *SpecsImage `json:"image"`
	GradientColors  []string    `json:"gradient_colors"`
}

type SpecsImage struct {
	Md5             *string `json:"md5"`
	Width           *uint32 `json:"width"`
	Height          *uint32 `json:"height"`
	BackgroundColor *string `json:"background_color"`
	TopOffset       *int32  `json:"top_offset"`
}

type SellerFlag struct {
	Name      *string `json:"name"`
	ImageFlag *Image  `json:"image_flag"`
}

type ImageOverlay struct {
	FeaturedPromoOverlay *FeaturedPromoOverlay `json:"featured_promo_overlay"`
	LegacyOverlay        *LegacyOverlay        `json:"legacy_overlay"`
}

type LegacyOverlay struct {
	Name         *string `json:"name"`
	ImageOverlay *Image  `json:"image_overlay"`
}

type FeaturedPromoOverlay struct {
	ItemId        *uint64                    `json:"item_id"`
	ModelId       *uint64                    `json:"model_id"`
	PromotionType *uint32                    `json:"promotion_type"`
	PromotionId   *uint64                    `json:"promotion_id"`
	OverlayType   *uint32                    `json:"overlay_type"`
	IsTeaser      *bool                      `json:"is_teaser"`
	AssetType     *uint32                    `json:"asset_type"`
	BasicOverlay  *BasicFeaturedPromoOverlay `json:"basic_overlay"`
	ImageOverlay  *Image                     `json:"image_overlay"`
}

type BasicFeaturedPromoOverlay struct {
	Icon            *Image                               `json:"icon"`
	Text            *string                              `json:"text"`
	TextColor       *string                              `json:"text_color"`
	BackgroundColor *FeaturedPromoOverlayBackgroundColor `json:"background_color"`
}

type FeaturedPromoOverlayBackgroundColor struct {
	IsGradient     *bool   `json:"is_gradient"`
	Color          *string `json:"color"`
	MainColor      *string `json:"main_color"`
	SecondaryColor *string `json:"secondary_color"`
}

type IconInImage struct {
	IconType *uint32 `json:"icon_type"`
	AdsText  *string `json:"ads_text"`
}

type ItemCardMask struct {
	MaskType *uint32 `json:"mask_type"`
	MaskText *string `json:"mask_text"`
}

type VideoInfo struct {
	VideoId       *string        `json:"video_id"`
	ThumbUrl      *string        `json:"thumb_url"`
	Duration      *int32         `json:"duration"`
	Version       *uint32        `json:"version"`
	Vid           *string        `json:"vid"`
	Formats       []*VideoFormat `json:"formats"`
	DefaultFormat *VideoFormat   `json:"default_format"`
	MmsData       *string        `json:"mms_data"`
}

type VideoFormat struct {
	Format  *uint32 `json:"format"`
	Defn    *string `json:"defn"`
	Profile *string `json:"profile"`
	Path    *string `json:"path"`
	Url     *string `json:"url"`
	Width   *uint32 `json:"width"`
	Height  *uint32 `json:"height"`
}

type Image struct {
	Hash   *string `json:"hash"`
	Width  *uint32 `json:"width"`
	Height *uint32 `json:"height"`
}

type ItemCardSpecs struct {
	CardBackgroundColor *CardBackgroundColor `json:"card_background_color"`
}

type CardBackgroundColor struct {
	Style          *string              `json:"style"`
	GradientColors []*CardGradientColor `json:"gradient_colors"`
}

type CardGradientColor struct {
	Color         *string  `json:"color"`
	OffsetPercent *float64 `json:"offset_percent"`
}

type DisplayPrice struct {
	Price                  *int64  `json:"price"`
	StrikethroughPrice     *int64  `json:"strikethrough_price"`
	HiddenPriceDisplayText *string `json:"hidden_price_display_text"`
}

type CTAButton struct {
	CtaType *uint32 `json:"cta_type"`
}
