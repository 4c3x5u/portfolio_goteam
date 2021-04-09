import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { Row } from 'react-bootstrap';
import { DragDropContext } from 'react-beautiful-dnd';

import Column from './Column/Column';
import { columnOrder } from './Column/columnOrder';

import './board.sass';

const Board = ({ handleActivate }) => {
  const [columns, setColumns] = useState([
    {
      id: 0,
      order: columnOrder.INBOX,
      tasks: [
        {
          id: 0,
          title: 'Some task',
          description: 'Some task description',
          order: 0,
        },
        {
          id: 1,
          title: 'Some other task',
          description: 'Some other task description',
          order: 1,
        },
      ],
    },
    { id: 1, order: columnOrder.READY, tasks: [] },
    { id: 2, order: columnOrder.GO, tasks: [] },
    { id: 3, order: columnOrder.DONE, tasks: [] },
  ]);

  // TODO: API call to the database here inside a useEffect
  useEffect(() => (
    console.log('Getting columns for board as well as tasks')
  ), []);

  const handleOnDragEnd = (result) => {
    if (!result.destination) return;

    const source = columns.find(
      (column) => column.id.toString() === result.source.droppableId,
    );

    const [item] = source.tasks.splice(result.source.index, 1);

    const newSource = {
      ...source,
      tasks: source.tasks.map(
        (task, index) => ({ ...task, order: index }),
      ),
    };

    const destination = columns.find(
      (column) => column.id.toString() === result.destination.droppableId,
    );

    destination.tasks.splice(result.destination.index, 0, item);

    const newDestination = {
      ...destination,
      tasks: destination.tasks.map(
        (task, index) => ({ ...task, order: index }),
      ),
    };

    const newColumns = (
      columns.map((column) => {
        switch (column.id) {
          case destination.id: return newDestination;
          case source.id: return newSource;
          default: return column;
        }
      })
    );

    setColumns(newColumns);
  };

  return (
    <Row id="Board">
      <DragDropContext onDragEnd={handleOnDragEnd}>
        {columns.map((column) => (
          <Column
            id={column.id}
            name={column.order}
            tasks={column.tasks}
            handleActivate={
              column.order === columnOrder.INBOX
                && handleActivate
            }
          />
        ))}
      </DragDropContext>
    </Row>
  );
};

Board.propTypes = { handleActivate: PropTypes.func.isRequired };

export default Board;
