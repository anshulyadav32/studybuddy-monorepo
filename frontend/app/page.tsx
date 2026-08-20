const courses = [
  { title: 'Mathematics Essentials', progress: 72, next: 'Algebraic identities', teacher: 'Aanya Sharma' },
  { title: 'Physics Foundations', progress: 54, next: 'Laws of motion', teacher: 'Rohan Iyer' },
  { title: 'Chemistry Lab', progress: 86, next: 'Acid-base titration', teacher: 'Meera Nair' },
];

const quickActions = ['Resume lesson', 'Practice quiz', 'View report card', 'Join live class'];
const highlights = [
  { label: 'Streak', value: '12 days' },
  { label: 'Hours learned', value: '18.4h' },
  { label: 'Avg. score', value: '88%' },
  { label: 'Focus score', value: '91/100' },
];

export default function HomePage() {
  return (
    <main style={{ padding: '40px 20px 80px', fontFamily: 'sans-serif', background: '#f8fafc' }}>
      <section style={{ maxWidth: 1200, margin: '0 auto' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: 16, marginBottom: 32, flexWrap: 'wrap' }}>
          <div>
            <p style={{ margin: 0, fontSize: 14, letterSpacing: '0.08em', color: '#475569', textTransform: 'uppercase' }}>Student dashboard</p>
            <h1 style={{ margin: '10px 0 0', fontSize: 38, lineHeight: 1.1 }}>Welcome back, Aisha</h1>
          </div>
          <button style={{ background: '#2563eb', color: '#fff', border: 'none', borderRadius: 12, padding: '12px 18px', fontSize: 16, fontWeight: 700 }}>Continue learning</button>
        </div>

        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(180px, 1fr))', gap: 18, marginBottom: 32 }}>
          {highlights.map((item) => (
            <div key={item.label} style={{ background: '#fff', borderRadius: 18, padding: '20px 18px', boxShadow: '0 10px 30px rgba(15, 23, 42, 0.06)', border: '1px solid #e2e8f0' }}>
              <div style={{ fontSize: 13, color: '#64748b', marginBottom: 12 }}>{item.label}</div>
              <div style={{ fontSize: 28, fontWeight: 800 }}>{item.value}</div>
            </div>
          ))}
        </div>

        <div style={{ display: 'grid', gridTemplateColumns: '1.6fr 0.9fr', gap: 22 }}>
          <div style={{ background: '#fff', borderRadius: 20, padding: 24, boxShadow: '0 10px 30px rgba(15, 23, 42, 0.06)', border: '1px solid #e2e8f0' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 20 }}>
              <h2 style={{ margin: 0, fontSize: 24 }}>Current courses</h2>
              <span style={{ color: '#2563eb', fontWeight: 700 }}>3 active</span>
            </div>

            <div style={{ display: 'grid', gap: 20 }}>
              {courses.map((course) => (
                <div key={course.title} style={{ background: '#f8fafc', borderRadius: 16, padding: 18, border: '1px solid #e2e8f0' }}>
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: 12, marginBottom: 12 }}>
                    <div>
                      <div style={{ fontSize: 18, fontWeight: 700 }}>{course.title}</div>
                      <div style={{ color: '#64748b', fontSize: 13, marginTop: 4 }}>with {course.teacher}</div>
                    </div>
                    <div style={{ fontWeight: 700, color: '#0f172a' }}>{course.progress}%</div>
                  </div>

                  <div style={{ background: '#dbeafe', borderRadius: 999, height: 10, overflow: 'hidden' }}>
                    <div style={{ background: '#2563eb', width: `${course.progress}%`, height: '100%' }} />
                  </div>

                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginTop: 12 }}>
                    <span style={{ fontSize: 14, color: '#475569' }}>Next: {course.next}</span>
                    <button style={{ background: '#e2e8f0', color: '#0f172a', border: 'none', borderRadius: 999, padding: '8px 12px', fontWeight: 700 }}>Resume</button>
                  </div>
                </div>
              ))}
            </div>
          </div>

          <div style={{ display: 'grid', gap: 22 }}>
            <div style={{ background: '#fff', borderRadius: 20, padding: 24, boxShadow: '0 10px 30px rgba(15, 23, 42, 0.06)', border: '1px solid #e2e8f0' }}>
              <h3 style={{ marginTop: 0, marginBottom: 18, fontSize: 20 }}>Quick actions</h3>
              <div style={{ display: 'grid', gap: 10 }}>
                {quickActions.map((action) => (
                  <button key={action} style={{ textAlign: 'left', background: '#f8fafc', border: '1px solid #e2e8f0', borderRadius: 12, padding: '12px 14px', fontSize: 15, fontWeight: 600, color: '#0f172a' }}>{action}</button>
                ))}
              </div>
            </div>

            <div style={{ background: 'linear-gradient(135deg, #1d4ed8 0%, #2563eb 100%)', color: '#fff', borderRadius: 20, padding: 24, boxShadow: '0 10px 30px rgba(37, 99, 235, 0.3)' }}>
              <div style={{ fontSize: 13, opacity: 0.8, textTransform: 'uppercase', letterSpacing: '0.08em' }}>Performance</div>
              <div style={{ fontSize: 36, fontWeight: 800, marginTop: 10 }}>Top 10%</div>
              <div style={{ marginTop: 8, lineHeight: 1.5, opacity: 0.9 }}>You are outperforming peers this week. Keep your momentum with 3 more practice sessions.</div>
            </div>
          </div>
        </div>
      </section>
    </main>
  );
}
