type LoadingDotsProps = {
  message?: string;
};

export default function LoadingDots({ message }: LoadingDotsProps) {
  return (
    <div className="feed-loading">
      <div className="dots">
        <span></span>
        <span></span>
        <span></span>
      </div>
      {message && <p className="feed-message">{message}</p>}
    </div>
  );
}
