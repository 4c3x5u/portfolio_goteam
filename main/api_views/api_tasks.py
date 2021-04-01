from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
from rest_framework.parsers import JSONParser
from ..models import Column, Task, Subtask
from ..serializers.ser_task import TaskSerializer
from ..serializers.ser_subtask import SubtaskSerializer
import json


@api_view(['POST'])
def tasks(request):
    column_tasks = Column.objects.filter(id=request.data.get('column'))
    task_serializer = TaskSerializer(
        data={'title': request.data.get('title'),
              'description': request.data.get('description'),
              'order': len(column_tasks) + 1,
              'column': request.data.get('column')}
    )
    if not task_serializer.is_valid():
        return Response(task_serializer.errors, 400)
    task = task_serializer.save()

    for order, title in enumerate(request.data.get('subtasks'), start=0):
        subtask_serializer = SubtaskSerializer(
            data={'title': title.get('title'),
                  'order': order,
                  'task': task.id}
        )
        if not subtask_serializer.is_valid():
            task.delete()
            return Response(subtask_serializer.errors, 400)
        subtask_serializer.save()

    return Response({
        'msg': 'Task creation successful.',
        'task_id': task.id
    }, 201)

