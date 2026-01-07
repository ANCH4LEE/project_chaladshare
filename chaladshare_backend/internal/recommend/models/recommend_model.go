package models

// Seed = โพสต์ต้นทาง (จากไลก์ล่าสุด) ที่ใช้เป็นตัวตั้ง
type Seedpost struct {
	PostID int
	Label  string
	Vec    []float64
}

// Candidate = โพสต์ผู้สมัครที่เอาไปเทียบ similarity
// Vec ใส่ json:"-" เพื่อไม่ต้องส่งกลับ frontend
type Candidatepost struct {
	PostID      int    `json:"post_id"`
	AuthorID    int    `json:"author_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CoverURL    string `json:"cover_url"`
	Visibility  string `json:"visibility"`

	AuthorName string `json:"authorName"`
	AuthorImg  string `json:"authorImg"`
	Tags       string `json:"tags"`

	LikeCount int  `json:"like_count"`
	IsLiked   bool `json:"is_liked"`
	IsSaved   bool `json:"is_saved"`

	Vec []float64 `json:"-"`
}

// RecommendPost = รูปแบบที่ส่งให้หน้าเว็บ (ตรงกับ PostCard เธอ)
type RecommendPost struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Tags       string `json:"tags"`
	Img        string `json:"img"`
	AuthorName string `json:"authorName"`
	AuthorImg  string `json:"authorImg"`

	IsLiked   bool `json:"is_liked"`
	LikeCount int  `json:"like_count"`
	IsSaved   bool `json:"is_saved"`

	Score float64 `json:"score,omitempty"`
}
