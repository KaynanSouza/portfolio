import React, { useState, useEffect } from 'react';
import axios from 'axios';

const DashboardComponent = () => {
  const [projects, setProjects] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchProjects = async () => {
      try {
        const response = await axios.get('http://localhost:8080/api/projects');
        setProjects(response.data);
      } catch (err) {
        console.error('Error fetching projects:', err);
        setError(err.response ? `${err.response.status}: ${err.response.data.error}` : err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchProjects();
  }, []);

  if (loading) {
    return <div>Loading projects...</div>;
  }

  if (error) {
    return <div>Error: {error}</div>;
  }

  return (
    <div>
      <h1>GitLab Projects</h1>
      {projects.length > 0 ? (
        <ul>
          {projects.map((project) => (
            <li key={project.id}>
              <a href={project.web_url} target="_blank" rel="noopener noreferrer">
                {project.name_with_namespace}
              </a>
            </li>
          ))}
        </ul>
      ) : (
        <p>No projects found.</p>
      )}
    </div>
  );
};

export default DashboardComponent;
