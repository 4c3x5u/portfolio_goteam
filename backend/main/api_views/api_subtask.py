from rest_framework.decorators import api_view
from rest_framework.response import Response
from ..models import Subtask
from ..serializers.ser_subtask import SubtaskSerializer

@api_view(['PATCH'])
def subtasks(request):
    current_subtask = Subtask.objects.get(id=request.data.get('id'))
    serializer = SubtaskSerializer(current_subtask,
                                   data=request.data.get('data'),
                                   partial=True)
    if serializer.is_valid():
        new_subtask = serializer.save()
        return Response({
            'msg': 'Subtask update successful.',
            'id': new_subtask.id
        }, 200)