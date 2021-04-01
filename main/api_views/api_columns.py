from rest_framework.decorators import api_view

# Not going to be used by the user, but rather by the app in the background to
# rapidly create columns for each newly created board
@api_view(['POST'])
def create_column(request):
    s