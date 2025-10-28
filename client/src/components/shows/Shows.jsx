import Show from "../show/Show.jsx";

const Shows = ({ shows, updateShowReview, message }) => {
    return (
        <div className="streams-grid">
            {shows && shows.length > 0
                ? shows.map((show) => (
                    <Show key={show._id} updateShowReview={updateShowReview} show={show} />
                ))
                : <h2>{message}</h2>
            }
        </div>
    )
}
export default Shows;