import sqlite3
import time
import logging
import os
from generate import ContentGenerator

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)

class ContentWorker:
    def __init__(self, db_path: str = "../db/content.db"):
        self.db_path = db_path
        self.generator = ContentGenerator()
        self.running = True
        
        # Ensure database directory exists
        os.makedirs(os.path.dirname(db_path), exist_ok=True)
        
        logging.info("Content worker initialized")

    def get_pending_job(self):
        """Get the next pending job from database"""
        try:
            conn = sqlite3.connect(self.db_path)
            cursor = conn.cursor()
            
            cursor.execute("""
                SELECT id, topic FROM jobs 
                WHERE status = 'pending' 
                ORDER BY created_at ASC 
                LIMIT 1
            """)
            
            job = cursor.fetchone()
            conn.close()
            
            if job:
                return {"id": job[0], "topic": job[1]}
            return None
            
        except sqlite3.Error as e:
            logging.error(f"Database error getting pending job: {e}")
            return None

    def update_job_status(self, job_id: int, status: str, output: str = ""):
        """Update job status and output in database"""
        try:
            conn = sqlite3.connect(self.db_path)
            cursor = conn.cursor()
            
            cursor.execute("""
                UPDATE jobs 
                SET status = ?, output = ?, updated_at = CURRENT_TIMESTAMP 
                WHERE id = ?
            """, (status, output, job_id))
            
            conn.commit()
            conn.close()
            
            logging.info(f"Job {job_id} updated to status: {status}")
            return True
            
        except sqlite3.Error as e:
            logging.error(f"Database error updating job {job_id}: {e}")
            return False

    def process_job(self, job):
        """Process a single job"""
        job_id = job["id"]
        topic = job["topic"]
        
        logging.info(f"Processing job {job_id}: {topic}")
        
        # Update status to processing
        if not self.update_job_status(job_id, "processing"):
            return
        
        try:
            # Generate content
            content = self.generator.generate_content(topic)
            
            if content:
                # Update with completed status and output
                self.update_job_status(job_id, "completed", content)
                logging.info(f"Job {job_id} completed successfully")
            else:
                self.update_job_status(job_id, "failed", "Content generation failed")
                logging.error(f"Job {job_id} failed: No content generated")
                
        except Exception as e:
            error_msg = f"Processing error: {str(e)}"
            self.update_job_status(job_id, "failed", error_msg)
            logging.error(f"Job {job_id} failed: {e}")

    def run(self):
        """Main worker loop"""
        logging.info("Starting content worker...")
        
        while self.running:
            try:
                # Check for pending jobs
                job = self.get_pending_job()
                
                if job:
                    self.process_job(job)
                else:
                    # No jobs available, wait before checking again
                    time.sleep(5)
                    
            except KeyboardInterrupt:
                logging.info("Worker stopped by user")
                self.running = False
                break
            except Exception as e:
                logging.error(f"Unexpected error in worker loop: {e}")
                time.sleep(10)  # Wait longer on unexpected errors

        logging.info("Content worker stopped")

def main():
    worker = ContentWorker()
    try:
        worker.run()
    except KeyboardInterrupt:
        logging.info("Shutting down worker...")

if __name__ == "__main__":
    main()
