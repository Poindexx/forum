package forum

type Session struct {
	ID            string
	Username      string
	Authenticated bool
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type ViewData struct {
	Title     string
	Username1 string
	Id1       string
	Posts     []Post1
	Post      Post1
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Likes struct {
	Id       int `json:"id"`
	Like     int `json:"like"`
	DisLike  int `json:"dislike"`
	AuthorID int `json:"author_id"`
	PostID   int `json:"psot_id"`
}

type LikeUpdateRequest struct {
	PostID   string `json:"postId"`
	Type     string `json:"type"`
	AuthorID int    `json:"authorId"`
}

type SortedRequest struct {
	Category_ids   []string `json:"category_ids"`
	Category_names []string `json:"category_names"`
	Sort_post      string   `json:"sort_post"`
	Text_dis       string   `json:"text_dis"`
	Start_date     string   `json:"start_date"`
	End_date       string   `json:"end_date"`
}

type PostDelRequest struct {
	Username string `json:"name"`
	PostId   string `json:"postid"`
}

type LikeUpdateComRequest struct {
	ComID    string `json:"comId"`
	Type     string `json:"type"`
	AuthorID int    `json:"authorId"`
}
type CommentRequest struct {
	PostID string `json:"post_id"`
}

type LikeUpdateResponse struct {
	NewLikeCount    int    `json:"newLikeCount"`
	NewDislikeCount int    `json:"newDislikeCount"`
	NewPostid       string `json:"newPostid"`
}
type LikeComUpdateResponse struct {
	NewLikeCount    int `json:"newLikeCount"`
	NewDislikeCount int `json:"newDislikeCount"`
	NewComid        int `json:"newComid"`
}
type CommResponse struct {
	Com string `json:"com"`
}

type LikesComments struct {
	Id       int      `json:"id"`
	Like     int      `json:"like"`
	DisLike  int      `json:"dislike"`
	AuthorID int      `json:"author_id"`
	Comments Comments `json:"comment_id"`
}

type Comments struct {
	Id       int    `json:"id"`
	PostID   string `json:"post_id"`
	Comment  string `json:"comment"`
	AuthorID int    `json:"author_id"`
	Author   string `json:"author"`
	Date     string `json:"date"`
}

type Comments1 struct {
	Id       int    `json:"id"`
	PostID   string `json:"post_id"`
	Comment  string `json:"comment"`
	AuthorID int    `json:"author_id"`
	Author   string `json:"author"`
	Date     string `json:"date"`
	Like     int    `json:"like"`
	DisLike  int    `json:"dislike"`
}

type Post struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Anons       string `json:"anons"`
	AuthorID    int    `json:"author_id"`
	Author      string `json:"author"`
	Date        string `json:"date"`
	ImageBase64 string `json:"imageBase64"`
	ImageName   string `json:"imageName"`
	CategoryIDs []string
	Categorys   []string
}

type Post1 struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Anons       string   `json:"anons"`
	AuthorID    int      `json:"author_id"`
	Author      string   `json:"author"`
	Date        string   `json:"date"`
	ImageURL    string   `json:"imageurl"`
	CategoryIDs []string `json:"category_id"`
	Categorys   []string `json:"category"`
	Like        int      `json:"like"`
	DisLike     int      `json:"doslike"`
	CommentLen  int      `json:"comment_len"`
	PostID      int      `json:"post_id"`
}
