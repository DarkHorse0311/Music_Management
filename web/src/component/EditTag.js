import { useState, useEffect } from "react";
import { useParams } from "react-router";

function EditTag() {
  let params = useParams();
  const [tag, setTag] = useState({});

  useEffect(() => {
    fetch("/api/v1/get_tag_info", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        id: parseInt(params.id),
      }),
    })
      .then((res) => res.json())
      .then((data) => {
        if (data.error) {
          alert(data.error);
        } else {
          setTag(data.tag);
        }
      });
  }, []);

  return (
    <div className="page">
      <h3>Edit Tag</h3>
      <div>
        <label htmlFor="id">ID</label>
        <input
          type="text"
          disabled
          name="id"
          id="id"
          value={tag.id}
          onChange={(e) => setTag({ ...tag, id: e.target.value })}
        />
        <label htmlFor="name">Name</label>
        <input
          type="text"
          name="name"
          id="name"
          value={tag.name}
          onChange={(e) => setTag({ ...tag, name: e.target.value })}
        />
        <label htmlFor="description">Description</label>
        <textarea
          name="description"
          id="description"
          value={tag.description}
          onChange={(e) => setTag({ ...tag, description: e.target.value })}
        />
        <button onClick={() => {}}>Save</button>
      </div>
    </div>
  );
}

export default EditTag;
