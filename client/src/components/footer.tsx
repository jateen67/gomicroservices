export default function Footer() {
  return (
    <div className="container">
      <div className="row">
        <div className="col text-center">
          <hr></hr>
          <small className="text-muted">
            &copy; {new Date().getFullYear()} Jatin Kalsi
          </small>
        </div>
      </div>
    </div>
  );
}
