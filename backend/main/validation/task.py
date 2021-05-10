from rest_framework.exceptions import ErrorDetail
from rest_framework.response import Response
from ..models import Task


# return (task, response)
def validate_task_id(task_id):
    if not task_id:
        return Response({
            'task_id': ErrorDetail(string='Task ID cannot be empty.',
                                   code='blank')
        }, 400)

    try:
        int(task_id)
    except ValueError:
        return Response({
            'task_id': ErrorDetail(string='Task ID must be a number.',
                                   code='invalid')
        }, 400)


