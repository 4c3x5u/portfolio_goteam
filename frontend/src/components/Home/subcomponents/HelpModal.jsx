import React from 'react';

import './helpmodal.sass';

const HelpModal = () => (
  <div id="HelpModal">
    <h1>Summary</h1>
    <p>
      In both the dropdowns and the inboz column, you can click
      <br />
      the plus icon to add a new item, and right click on items to view
      additional controls for that specific item.
    </p>
    <h1>Team Controls</h1>
    <ol>
      <li>
        Click the TEAM icon located on the navbar to toggle the Team drop-down
        menu.
      </li>
      <li>
        The team members that are included in the currently displayed board,
        and therefore can see and interact with items, will have a yellow tick
        on both sides of their name.
      </li>
      <li>
        Click a member to include or disclude him from the current board.
      </li>
      <li>
        You can click the plus icon at the bottom of the list to invite new
        members to your team.
      </li>
    </ol>
    <h1>Boards Controls</h1>
    <ol>
      <li>
        Click the BOARDS icon located on the navbar to toggle the Boards
        drop-down menu.
      </li>
      <li>
        The currently displayed board will have arrows pointing to its name
        on both of its name.
      </li>
      <li>Click a board to view its content.</li>
      <li>
        You can click the plus icon at the bottom of the list to create a new
        board.
      </li>
    </ol>
    <h1>Task Controls</h1>
    <ol>
      <li>
        Click the plus icon inside the INBOX column to create a new task.
      </li>
      <li>
        Give the task a name. You can optionally give it a description and/or
        associated subtasks.
      </li>
      <li>
        Drag and drop tasks from one column to another in order to move them
        around.
      </li>
      <li>
        Right click tasks to edit, delete, or assign team members to them.
      </li>
    </ol>
  </div>
);

export default HelpModal;
