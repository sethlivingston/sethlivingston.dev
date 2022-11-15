import { remarkDefinitionList, defListHastHandlers } from "remark-definition-list";
import { unified } from "unified";
import remarkParse from "remark-parse";
import remarkRehype from "remark-rehype";
import rehypeStringify from "rehype-stringify";

const md = `
Test for defList.

Apple
:   Pomaceous fruit of plants of the genus Malus in
    the family Rosaceae.

Orange
:   The fruit of an evergreen tree of the genus Citrus.
`;

const html = await unified()
  .use(remarkParse)
  .use(remarkDefinitionList)
  .use(remarkRehype, {
    handlers: {
      // any other handlers
      ...defListHastHandlers,
    },
  })
  .use(rehypeStringify)
  .process(md);

console.log(html);
