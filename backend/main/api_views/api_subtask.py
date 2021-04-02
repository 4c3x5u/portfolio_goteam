from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
from ..models import Subtask
from ..serializers.ser_subtask import SubtaskSerializer


@api_view(['PATCH'])
def subtasks(request):
    subtask_id = request.data.get('id')
    if not subtask_id:
        return Response({
            'id': ErrorDetail(string='Subtask ID cannot be empty.',
                              code='blank')
        }, 400)

    serializer = SubtaskSerializer(Subtask.objects.get(id=subtask_id),
                                   data=request.data.get('data'),
                                   partial=True)
    if serializer.is_valid():
        new_subtask = serializer.save()
        return Response({
            'msg': 'Subtask update successful.',
            'id': new_subtask.id
        }, 200)
