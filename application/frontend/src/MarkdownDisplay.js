import React, { useEffect, useState } from 'react';
import marked from 'marked';

const MarkdownDisplay = () => {
  const [markdown, setMarkdown] = useState("Loading...");

  useEffect(() => {
    fetch('/yourfile.md')
      .then(response => response.text())
      .then(text => setMarkdown(marked(text)));
  }, []);

  return <div dangerouslySetInnerHTML={{ __html: markdown }} />;
};

export default MarkdownDisplay;
