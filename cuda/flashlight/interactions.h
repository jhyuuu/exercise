#ifndef __INTERACTIONS_H_
#define __INTERACTIONS_H_

#define W 600
#define H 600
#define DELTA 5
#define TITLE_STRING "flashlight"
int2 loc = {W/2, H/2};
bool dragMode = false;

void keyboard(unsigned char key, int x, int y) {
  if (key == 'a') dragMode = !dragMode;
  if (key == 27) exit(0); // ESC

  glutPostRedisplay();
}

void mouseMove(int x, int y) {
  if (dragMode) return;
  loc.x = x;
  loc.y = y;
  
  glutPostRedisplay();
}

void mouseDrag(int x, int y) {
  if (!dragMode) return;
  loc.x = x;
  loc.y = y;
  glutPostRedisplay();
}

void handleSpecialKeypress(int key, int x, int y) {
  if (key == GLUT_KEY_LEFT)  loc.x -= DELTA;
  if (key == GLUT_KEY_RIGHT) loc.x += DELTA;
  if (key == GLUT_KEY_UP)    loc.y -= DELTA;
  if (key == GLUT_KEY_DOWN)  loc.y += DELTA;
  glutPostRedisplay();
}

void printInstructions() {
  printf("flashlight interactions\n");
  printf("a: toggle mouse tracking mode\n");
  printf("arrow keys: move ref location\n");
  printf("esc: close graphics window\n");
}

#endif // __INTERACTIONS_H_
